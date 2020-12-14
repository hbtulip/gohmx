package grpcpool

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/keepalive"
)

var (
	ErrStringSplit    = errors.New("err string split")
	ErrNotFoundClient = errors.New("not found grpc conn")
	ErrConnShutdown   = errors.New("grpc conn shutdown")

	defaultClientPoolCap    = 2 //5
	defaultDialTimeout      = 5 * time.Second
	defaultKeepAlive        = 30 * time.Second
	defaultKeepAliveTimeout = 10 * time.Second
)

type ClientOption struct {
	DialTimeout      time.Duration
	KeepAlive        time.Duration
	KeepAliveTimeout time.Duration
	ClientPoolSize   int
}

func NewDefaultClientOption() *ClientOption {
	return &ClientOption{
		DialTimeout:      defaultDialTimeout,
		KeepAlive:        defaultKeepAlive,
		KeepAliveTimeout: defaultKeepAliveTimeout,
	}
}

type ClientPool struct {
	option   *ClientOption
	capacity int64
	next     int64
	target   string

	sync.Mutex

	conns []*grpc.ClientConn
}

func NewClient(target string, option *ClientOption) *ClientPool {
	if option == nil {
		option = NewDefaultClientOption()
	}
	if option.ClientPoolSize <= 0 {
		option.ClientPoolSize = defaultClientPoolCap
	}

	var cc = ClientPool{
		target:   target,
		conns:    make([]*grpc.ClientConn, option.ClientPoolSize),
		capacity: int64(option.ClientPoolSize),
		option:   option,
		next:     int64(0),
	}

	for idx, _ := range cc.conns {
		conn, _ := cc.connect()
		cc.conns[idx] = conn
	}

	return &cc
}

func (cc *ClientPool) GetClientConn() (*grpc.ClientConn, error) {
	var (
		idx  int64
		next int64

		err error
	)

	next = atomic.AddInt64(&cc.next, 1)
	idx = next % cc.capacity
	conn := cc.conns[idx]
	if conn != nil && cc.checkState(conn) == nil {
		return conn, nil
	}

	// gc old conn
	if conn != nil {
		conn.Close()
	}

	cc.Lock()
	defer cc.Unlock()

	// double check, already inited
	conn = cc.conns[idx]
	if conn != nil && cc.checkState(conn) == nil {
		return conn, nil
	}

	conn, err = cc.connect()
	if err != nil {
		return nil, err
	}

	cc.conns[idx] = conn
	return conn, nil
}

func (cc *ClientPool) checkState(conn *grpc.ClientConn) error {
	state := conn.GetState()
	switch state {
	case connectivity.TransientFailure, connectivity.Shutdown:
		return ErrConnShutdown
	}

	return nil
}

func (cc *ClientPool) connect() (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(cc.target,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(cc.option.DialTimeout),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:    cc.option.KeepAlive,
			Timeout: cc.option.KeepAliveTimeout},
		),
	)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (cc *ClientPool) Close() {
	cc.Lock()
	defer cc.Unlock()

	for _, conn := range cc.conns {
		if conn == nil {
			continue
		}

		conn.Close()
	}
}
