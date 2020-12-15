package main

import (
	"gohmx/controller"
	"gohmx/model"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var exitChan chan os.Signal

func exitHandle() {

	s := <-exitChan
	log.Println("收到退出信号：", s)

	model.CloseGrpcPool()

	os.Exit(1)
}

func main() {

	exitChan = make(chan os.Signal)
	//syscall.SIGINT 2,syscall.SIGKILL 9,syscall.SIGTERM 15
	signal.Notify(exitChan, os.Interrupt, os.Kill, syscall.SIGTERM)
	go exitHandle()

	controller.InitRouter()

	log.Println("Gohmx exit!")
}
