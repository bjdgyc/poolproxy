package main

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/naoina/toml"
)

// 配置文件
type Config struct {
	Logfile string
	Options map[string]Option
}

type Option struct {
	// 日志文件
	Logfile string
	// ip:port 或者socket文件路径
	Addr string `toml:"addr"`

	// 最大等待时间 默认是5秒
	PoolTimeout time.Duration `toml:"pool_timeout"`
	// 读取超时时间
	ReadTimeout time.Duration `toml:"read_timeout"`
	// 写入超时时间
	WriteTimeout time.Duration `toml:"write_timeout"`

	// 远程配置
	RAddr string `toml:"raddr"`
	RUser string `toml:"ruser"`
	RPass string `toml:"rpass"`
	// 远程最大连接数 默认是10秒
	RPoolSize int `toml:"rpool_size"`
	// 远程KeepAlive间隔时间
	RKeepAlivePeriod time.Duration `toml:"rkeep_alive_period"`
	// 远程最大空闲时间 默认是2分钟 120
	RIdleTimeout time.Duration `toml:"ridle_timeout"`
	// 远程空闲连接检测 默认3分钟一次 180
	RIdleCheckFrequency time.Duration `toml:"ridle_check_frequency"`
}

func (opt *Option) init() error {

	if opt.RPoolSize == 0 {
		opt.RPoolSize = 2
	}

	// 时间参数转换为秒
	opt.ReadTimeout = opt.ReadTimeout * time.Second
	opt.WriteTimeout = opt.WriteTimeout * time.Second
	opt.PoolTimeout = opt.PoolTimeout * time.Second

	opt.RKeepAlivePeriod = opt.RKeepAlivePeriod * time.Second
	opt.RIdleTimeout = opt.RIdleTimeout * time.Second
	opt.RIdleCheckFrequency = opt.RIdleCheckFrequency * time.Second

	// 设置默认值
	if opt.RKeepAlivePeriod == 0 {
		opt.RKeepAlivePeriod = 5 * time.Second
	}
	if opt.PoolTimeout == 0 {
		opt.PoolTimeout = 5 * time.Second
	}
	if opt.RIdleTimeout == 0 {
		opt.RIdleTimeout = time.Second * 120
	}
	if opt.RIdleCheckFrequency == 0 {
		opt.RIdleCheckFrequency = time.Second * 180
	}

	return nil
}

func LoadConfig(file string) *Config {
	f, err := os.Open(file)
	if err != nil {
		commonLog.Fatal(err)
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		commonLog.Fatal(err)
	}
	config := new(Config)
	err = toml.Unmarshal(buf, config)
	if err != nil {
		commonLog.Fatal(err)
	}
	return config
}
