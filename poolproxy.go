package main

import (
	"flag"
	"fmt"
	"github.com/bjdgyc/slog"
	"net/http"
	_ "net/http/pprof"
)

var (
	connPool  *ConnPool
	commonLog *slog.Logger = slog.GetStdLog()
)

func main() {

	cfile := flag.String("c", "./config.toml", "配置文件")
	flag.Parse()

	config := LoadConfig(*cfile)

	if config.Logfile != "" {
		commonLog = slog.New(config.Logfile, "")
	}

	for _, opt := range config.Options {

		logger := commonLog
		if opt.Logfile != "" {
			logger = slog.New(opt.Logfile, "")
		}

		connPool = NewConnPool(opt, logger)
		go StartProxy(connPool, opt.Addr)
		fmt.Println(connPool)
	}

	go func() {
		http.ListenAndServe("127.0.0.1:8090", nil)
	}()

	//connPool.Close()
	InitSignal()

}
