package server

import (
	"flow/config"
	"flow/traffic"
	"flow/utils"
	"fmt"
	"go.uber.org/zap"
	"net"
	"net/http"
	"net/http/httputil"
)

func StartTCPServer() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.XConfig.Server.TCPPort))
	if err != nil {
		zap.L().Fatal("Error starting TCP server", zap.Error(err))
	}
	defer listener.Close()

	zap.L().Info("Starting TCP server", zap.Int("port", config.XConfig.Server.TCPPort))

	for {
		conn, err := listener.Accept()
		if err != nil {
			zap.L().Error("Error accepting connection", zap.Error(err))
			continue
		}
		go utils.HandleTCP(conn)
	}
}

func StartUDPServer() {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", config.XConfig.Server.UDPPort))
	if err != nil {
		zap.L().Fatal("Error starting UDP server", zap.Error(err))
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		zap.L().Fatal("Error starting UDP server", zap.Error(err))
	}
	defer conn.Close()

	zap.L().Info("Starting UDP server", zap.Int("port", config.XConfig.Server.UDPPort))

	buf := make([]byte, 1024)
	for {
		n, clientAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			zap.L().Error("Error reading from UDP client", zap.Error(err))
			continue
		}
		go utils.HandleUDP(conn, clientAddr, buf[:n])
	}
}

func StartHttpOrHttpsServer() {
	director := func(req *http.Request) {
		// 设置目标服务器地址
		req.URL.Scheme = "http"
		req.URL.Host = config.XConfig.WebServer.RemoteAddr
	}

	proxy := &httputil.ReverseProxy{Director: director}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		rAddr, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			zap.L().Error("Error parsing IP", zap.Error(err))
			return
		}

		traffic.RecordRequest(rAddr)
		traffic.Tacker <- struct{}{}
		if traffic.IsBlocked(rAddr) {
			zap.L().With(zap.String("TAG", "Blocked")).Warn("UDP packet blocked", zap.String("remote_addr", rAddr))
			// 尝试接管连接
			hijacker, ok := w.(http.Hijacker)
			if !ok {
				http.Error(w, "HTTP/1.1 required to hijack", http.StatusInternalServerError)
				return
			}

			conn, _, err := hijacker.Hijack()
			if err != nil {
				http.Error(w, "Failed to hijack", http.StatusInternalServerError)
				return
			}
			// 关闭连接
			conn.Close()
			return
		}

		proxy.ServeHTTP(w, r)
	})

	if config.XConfig.WebServer.EnableTLS {
		zap.L().Info("Starting HTTPS server", zap.String("host_addr", config.XConfig.WebServer.HostAddr))
		// 配置 HTTPS 服务器
		err := http.ListenAndServeTLS(config.XConfig.WebServer.HostAddr, config.XConfig.WebServer.CertFile, config.XConfig.WebServer.KeyFile, nil)
		if err != nil {
			zap.L().Fatal("Error starting HTTPS server", zap.Error(err))
			return
		}
	} else {
		zap.L().Info("Starting HTTP server", zap.String("host_addr", config.XConfig.WebServer.HostAddr))
		err := http.ListenAndServe(config.XConfig.WebServer.HostAddr, nil)
		if err != nil {
			zap.L().Fatal("Error starting HTTP server", zap.Error(err))
			return
		}
	}
}
