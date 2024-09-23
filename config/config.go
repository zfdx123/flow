package config

import (
	"flag"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Handle struct {
		RequestThreshold int           `yaml:"request_threshold"`
		BlockDuration    time.Duration `yaml:"block_duration"`
		RequestTimeout   time.Duration `yaml:"request_timeout"`
	} `yaml:"handle"`
	WebServer struct {
		EnableHttp bool   `yaml:"enable_http"`
		EnableTLS  bool   `yaml:"enable_tls"`
		HostAddr   string `yaml:"host_addr"`
		CertFile   string `yaml:"cert_file"`
		KeyFile    string `yaml:"key_file"`
		RemoteAddr string `yaml:"remote_addr"`
	} `yaml:"web_server"`
	Server struct {
		TCPPort         int    `yaml:"tcp_port"`
		UDPPort         int    `yaml:"udp_port"`
		EnableTCP       bool   `yaml:"enable_tcp"`
		InternalTCPAddr string `yaml:"internal_tcp_addr"`
		EnableUDP       bool   `yaml:"enable_udp"`
		InternalUDPAddr string `yaml:"internal_udp_addr"`
	} `yaml:"server"`
	Log struct {
		LevelFile    string `yaml:"level_file"`
		LevelConsole string `yaml:"level_console"`
		Path         string `yaml:"path"`
		MaxSize      int    `yaml:"max_size"`
		MaxAge       int    `yaml:"max_age"`
		MaxBackups   int    `yaml:"max_backups"`
		Compress     bool   `yaml:"compress"`
	} `yaml:"log"`
}

var XConfig Config

var configPath *string

func init() {
	configPath = flag.String("config", "config.yaml", "config path")
	flag.Parse()
}

// LoadConfig 从配置文件中读取配置
func LoadConfig() {
	file, err := os.Open(*configPath)
	if err != nil {
		log.Fatalf("Failed to open config file %v", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&XConfig); err != nil {
		log.Fatalf("Error decoding config file: %v", err)
	}
}
