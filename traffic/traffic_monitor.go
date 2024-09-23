package traffic

import (
	"flow/config"
	"go.uber.org/zap"
	"sync"
	"time"
)

var (
	blacklist   = make(map[string]time.Time)
	rateLimit   = make(map[string]int)
	lastRequest = make(map[string]time.Time)
	muBlacklist sync.RWMutex // 用于黑名单的读写锁
	muRateLimit sync.RWMutex // 用于流量监控的读写锁
)

// BlockIP 封禁IP
func BlockIP(ip string) {
	muBlacklist.Lock()
	defer muBlacklist.Unlock()
	blacklist[ip] = time.Now().Add(config.XConfig.Handle.BlockDuration)
	zap.L().With(zap.String("TAG", "Blocked")).Warn("Blocked IP", zap.String("ip", ip))
}

// IsBlocked 检查是否被封禁
func IsBlocked(ip string) bool {
	muBlacklist.RLock()
	defer muBlacklist.RUnlock()

	blockedUntil, found := blacklist[ip]
	if found {
		if time.Now().Before(blockedUntil) {
			zap.L().Info("Blocked IP", zap.String("ip", ip), zap.Time("blocked_until", blockedUntil))
			return true
		}
		delete(blacklist, ip) // 解除封禁
	}
	return false
}

// RecordRequest 记录请求
func RecordRequest(ip string) {
	zap.L().Info("Received request from", zap.String("remote_addr", ip))

	muRateLimit.Lock()
	defer muRateLimit.Unlock()
	// 更新请求计数和最后请求时间
	rateLimit[ip]++
	lastRequest[ip] = time.Now()
	zap.L().Info("Recorded request", zap.String("ip", ip), zap.Int("count", rateLimit[ip]))

	// 检查封禁状态
	if IsBlocked(ip) {
		zap.L().Info("Recorded request denied", zap.String("ip", ip))
		return
	}
}

var Tacker = make(chan struct{})

// MonitorTraffic 监控流量
func MonitorTraffic() {
	ticker := time.NewTicker(time.Second * 10) // 每 10 秒监控一次
	defer ticker.Stop()

	for {
		select {
		case <-Tacker:
			blackCheck()
		case <-ticker.C:
			blackCheck()
		}
	}
}

func blackCheck() {
	var toBlock []string

	muRateLimit.Lock()
	zap.L().Warn("Monitoring traffic... Current IP counts", zap.Any("rate_limit", rateLimit))
	zap.L().Warn("Blocked", zap.Any("black_limit", blacklist))
	for ip, count := range rateLimit {
		if count > config.XConfig.Handle.RequestThreshold {
			toBlock = append(toBlock, ip)
		}
	}

	// 检查请求计数的超时
	for ip, lastTime := range lastRequest {
		if time.Since(lastTime) > config.XConfig.Handle.RequestTimeout {
			delete(rateLimit, ip)   // 重置请求计数
			delete(lastRequest, ip) // 清除最后请求时间
			zap.L().Warn("Reset request count for IP", zap.String("ip", ip))
		}
	}

	muRateLimit.Unlock()

	// 封禁IP
	for _, ip := range toBlock {
		BlockIP(ip)
	}
}
