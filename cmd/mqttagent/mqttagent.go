package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"

	"gohmx/api"
)

/*
type OilTankServer interface {
	//获取油罐信息
	GetTankStatus(context.Context, *TankRequestParams) (*TankStatus, error)
}

// 油罐状态
message TankStatus {
	required  string mdc_id = 1;//box序列号
	required  uint32 comm_terminal_id = 2;//通讯终端编号
	required  uint32 tank_id = 3;//油罐ID
	required  uint32 probe_id = 4;//探棒ID
	required  uint64 datetime = 5;//时间戳,EPOCH秒数
	required  uint32 volume = 6;//油体积,单位0.01
	required  uint32 tc_volume = 7;//净油体积,单位0.01
	required  uint32 ullage = 8;//空容升数,单位0.01
	required  uint32 height = 9;//油面高,mm
	required  uint32 water_height = 10;//水面高,mm
	required  uint32 temperature = 11;//温度
	required  uint64 record_counter = 12;//计数器
}

*/

type server struct{}

func (s *server) GetTankStatus(ctx context.Context, a *OilMachine.TankRequestParams) (*OilMachine.TankStatus, error) {

	MdcId := "AX001"
	TankId := a.GetTankID()
	Height := uint32(1200)
	Temperature := uint32(18)
	Volume := uint32(36000)

	var res OilMachine.TankStatus
	res.MdcId = MdcId
	res.TankId = TankId
	res.Height = Height
	res.Temperature = Temperature
	res.Volume = Volume

	return &res, nil

}

func main() {

	lis, err := net.Listen("tcp", ":8010") //监听所有网卡8010端口的TCP连接
	if err != nil {
		log.Fatalf("监听失败: %v", err)
	}
	s := grpc.NewServer() //创建gRPC服务

	/**注册接口服务
	 * 以定义proto时的service为单位注册，服务中可以有多个方法
	 * (proto编译时会为每个service生成Register***Server方法)
	 * 包.注册服务方法(gRpc服务实例，包含接口方法的结构体[指针])
	 */
	//OilMachine.RegisterOilGunServer(s, &server{})
	OilMachine.RegisterOilTankServer(s, &server{})

	// 在gRPC服务器上注册反射服务
	reflection.Register(s)
	// 将监听交给gRPC服务处理
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
