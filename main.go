package main

import (
	"flow/config"
	"flow/log"
	"flow/server"
	"flow/traffic"
)

func main() {
	config.LoadConfig()         // 加载配置
	log.InitLogger()            // 初始化日志
	go traffic.MonitorTraffic() // 启动流量监控
	if config.XConfig.Server.EnableTCP {
		go server.StartTCPServer() // 启动TCP服务器
	}
	if config.XConfig.Server.EnableUDP {
		go server.StartUDPServer() // 启动UDP服务器
	}
	if config.XConfig.WebServer.EnableHttp {
		go server.StartHttpOrHttpsServer() // 启动HTTP或HTTPS服务器
	}

	select {} // 保持服务运行

}
