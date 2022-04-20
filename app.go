package main

import (
	"os"
	"os/exec"
	"os/signal"

	"github.com/wode-czw/trans_czw_wode/server"
	"github.com/wode-czw/trans_czw_wode/server/config"
)

func main() {
	//gin开一个携程比下面cmd开一个进程的速度会快得多
	go server.Run()

	start_Gin()

}

func start_Gin() {
	//启动chrome
	chSignal := make(chan os.Signal, 1) //这个东西等会用来接收系统的信号
	//signal.Notify(chSignal, syscall.SIGINT, syscall.SIGTERM)
	signal.Notify(chSignal, os.Interrupt)

	ChromePath := "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe"
	cmd := exec.Command(ChromePath, "--app=http://127.0.0.1"+config.Get_port()+"/static/index.html")

	cmd.Start()
	<-chSignal
	cmd.Process.Kill()
}
