package main

import (
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"poolproxy/slog"
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
		commonLog.SetLogfile(config.Logfile)
	}

	for _, opt := range config.Options {

		if opt.Logfile != "" {
			commonLog.SetLogfile(opt.Logfile)
		}

		connPool = NewConnPool(opt, commonLog)
		go StartProxy(connPool, opt.Addr)
		fmt.Println(opt.Addr, connPool)
	}

	go func() {
		http.ListenAndServe("127.0.0.1:8090", nil)
	}()

	// connPool.Close()
	InitSignal()

}
