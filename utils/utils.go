package utils

import (
	"flow/config"
	"flow/traffic"
	"go.uber.org/zap"
	"io"
	"net"
)

func HandleTCP(conn net.Conn) {
	defer conn.Close()
	ip := conn.RemoteAddr().String()

	rAddr, _, err := net.SplitHostPort(ip)
	if err != nil {
		zap.L().Error("Error parsing IP", zap.Error(err))
		return
	}
	traffic.RecordRequest(rAddr)
	traffic.Tacker <- struct{}{}
	if traffic.IsBlocked(rAddr) {
		zap.L().With(zap.String("TAG", "Blocked")).Warn("Connection blocked", zap.String("ip", ip))
		return
	}

	internalConn, err := net.Dial("tcp", config.XConfig.Server.InternalTCPAddr)
	if err != nil {
		zap.L().Error("Error connecting to internal server", zap.Error(err))
		return
	}
	defer internalConn.Close()

	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		if n == 0 || err == io.EOF {
			return
		}

		if err != nil {
			zap.L().Error("Error reading from client", zap.Error(err))
			return
		}
		_, err = internalConn.Write(buf[:n])
		if err != nil {
			zap.L().Error("Error writing to internal server", zap.Error(err))
			return
		}
	}
}

func HandleUDP(conn *net.UDPConn, addr *net.UDPAddr, data []byte) {
	ip := addr.String()

	rAddr, _, err := net.SplitHostPort(ip)
	if err != nil {
		zap.L().Error("Error parsing IP", zap.Error(err))
		return
	}

	traffic.RecordRequest(rAddr)
	traffic.Tacker <- struct{}{}
	if traffic.IsBlocked(rAddr) {
		zap.L().With(zap.String("TAG", "Blocked")).Warn("UDP packet blocked", zap.String("remote_addr", ip))
		return
	}

	internalAddr, err := net.ResolveUDPAddr("udp", config.XConfig.Server.InternalUDPAddr)
	if err != nil {
		zap.L().Error("Error resolving internal UDP address", zap.Error(err))
		return
	}

	internalConn, err := net.DialUDP("udp", nil, internalAddr)
	if err != nil {
		zap.L().Error("Error connecting to internal UDP server", zap.Error(err))
		return
	}
	defer internalConn.Close()

	_, err = internalConn.Write(data)
	if err != nil {
		zap.L().Error("Error forwarding data to internal server", zap.Error(err))
		return
	}

	// 动态读取数据
	buf := make([]byte, 0, 1024) // 初始容量
	for {
		tempBuf := make([]byte, 1024) // 临时缓冲区
		n, _, err := internalConn.ReadFromUDP(tempBuf)
		if err != nil {
			if n == 0 {
				break // 如果没有读取到数据，退出循环
			}
			zap.L().Error("Error reading from internal server", zap.Error(err))
			return
		}
		buf = append(buf, tempBuf[:n]...) // 动态扩展
	}

	_, err = conn.WriteToUDP(buf, addr)
	if err != nil {
		zap.L().Error("Error sending response to client", zap.Error(err))
	}
}
