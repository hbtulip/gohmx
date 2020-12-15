package model

import (
	"context"
	_ "google.golang.org/grpc"
	_ "google.golang.org/grpc/balancer/grpclb"
	"log"
	"strconv"

	"gohmx/api"
	"gohmx/common/grpcpool"
)

var oilClientPool *grpcpool.ClientPool

func init() {
	log.Println("Init gRPC ClientPool")
	oilClientPool = grpcpool.NewClient("localhost:8010", nil)

}

func CloseGrpcPool() {
	if oilClientPool != nil {
		log.Println("Close gRPC ClientPool")
		oilClientPool.Close()
	}
}

func GetTankStatus(storeid_s string, tankid_s string) (interface{}, error) {

	storeid, _ := strconv.ParseUint(storeid_s, 10, 64)
	tmp, _ := strconv.ParseUint(tankid_s, 10, 32)
	tankid := uint32(tmp)

	// 连接gRPC服务器，
	//conn, err := grpc.Dial("localhost:8010", grpc.WithInsecure())
	conn, err := oilClientPool.GetClientConn()

	/*
		// 指定IP轮询，负载均衡生产上可使用 consul，etcd 作服务发现
		conn, err := grpc.Dial("", grpc.WithInsecure(),
			grpc.WithBalancer(
				grpc.RoundRobin(
					grpclb.NewConsulResolver(
						"localhost:8010", "otherhost:8010",
					))),
		)
	*/
	if err != nil {
		return nil, err
	}
	//由客户连接池管理
	//defer conn.Close()

	// 连接GRPC
	grc := OilMachine.NewOilTankClient(conn)

	// 创建要发送的结构体
	req := OilMachine.TankRequestParams{
		StoreID: storeid,
		TankID:  tankid,
	}
	// 调用server的注册方法
	r, err := grc.GetTankStatus(context.Background(), &req)
	if err != nil {
		return nil, err
	}

	myMap := make(map[string]*OilMachine.TankStatus)
	myMap["TankStatus"] = r

	return myMap, nil
}
