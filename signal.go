package main

import (
	"os"
	"os/signal"
	"syscall"
)

var (
	exitChan = make(chan os.Signal, 1)
)

// InitSignal register signals handler.
func InitSignal() {
	signal.Notify(exitChan, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT) // , syscall.SIGSTOP
	for {
		s := <-exitChan
		commonLog.Info("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT: // , syscall.SIGSTOP
			return
		case syscall.SIGHUP:
			reload()
		default:
			return
		}
	}
}

// TODO 没有实现
func reload() {
}
