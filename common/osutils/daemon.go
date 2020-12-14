package osutils

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	//"io/ioutil"
	"runtime"
	//"strconv"
)

func doSomething() {
	done := make(chan struct{})
	sig := make(chan os.Signal)
	signal.Notify(sig)
	go func() {
		for s := range sig {
			if s == syscall.SIGQUIT || s == syscall.SIGINT {
				close(done)
				return
			}
		}
	}()

	for {
		time.Sleep(1 * time.Second)
		fmt.Println("do something")
		select {
		case <-done:
			fmt.Println("quit")
			return
		default:
		}
	}
}

func main() {
	dir, _ := os.Getwd()
	fmt.Println("当前路径：", dir)

	sysType := runtime.GOOS
	if sysType != "linux" {
		return
	}

	if len(os.Args) <= 1 {
		doSomething()
		return
	}
	/*
		cmd := os.Args[1]
		if cmd == "start" {
			pid, _ := syscall.ForkExec(os.Args[0], os.Args[:1], &syscall.ProcAttr{
				Env: append(os.Environ(), []string{"DAEMON=true"}...),
				Sys: &syscall.SysProcAttr{
					Setsid: true,
				},
				Files: []uintptr{0, 1, 2},
			})
			syscall.Umask(0000)
			_ = ioutil.WriteFile("/var/run/serviced.pid", []byte(strconv.Itoa(pid)), 0660)
			return
		} else if cmd == "stop" {
			reader, _ := os.Open("/var/run/serviced.pid")
			pidStr, _ := ioutil.ReadAll(reader)
			pid, _ := strconv.Atoi(string(pidStr))
			fmt.Println(pid)
			_ = syscall.Kill(pid, syscall.SIGQUIT)
			return
		}
	*/
}
