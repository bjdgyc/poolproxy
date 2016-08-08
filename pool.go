package main

import (
	"errors"
	"github.com/bjdgyc/slog"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrClosed      = errors.New("redis: coon pool is closed")
	ErrPoolTimeout = errors.New("redis: conn pool timeout")
	errConnActive  = errors.New("conn is not active")
)

var timers = sync.Pool{
	New: func() interface{} {
		return time.NewTimer(0)
	},
}

// PoolStats 连接池状态数据
type PoolStats struct {
	Requests  int32 // 获取连接池的次数
	Hits      int32 // 空闲连接池命中次数
	Timeouts  int32 // 获取连接池超时的次数
	UseConns  int32 // 使用中的连接池数据量
	FreeConns int32 // 空闲的连接池数据量
}

func (s *PoolStats) String() string {
	str := ""
	str += " Requests:" + strconv.Itoa(int(s.Requests))
	str += " Hits:" + strconv.Itoa(int(s.Hits))
	str += " Timeouts:" + strconv.Itoa(int(s.Timeouts))
	str += " UseConns:" + strconv.Itoa(int(s.UseConns))
	str += " FreeConns:" + strconv.Itoa(int(s.FreeConns))
	return str
}

type Pooler interface {
	Get() (*Conn, error)
	Put(*Conn, bool)
	FreeLen() int
	Stats() *PoolStats
	Close() error
}

type dialer func() (net.Conn, error)

type ConnPool struct {
	queue     chan struct{} // 连接池控制队列
	lock      sync.Mutex    // 并发锁
	freeConns []*Conn       // 空闲连接池
	stats     PoolStats
	opt       Option // 配置参数
	closed    bool   // 连接池关闭标志
	log       *slog.Logger
}

// 验证ConnPool是否实现了Pooler接口
var _ Pooler = (*ConnPool)(nil)

func NewConnPool(opt Option, logger *slog.Logger) *ConnPool {
	err := opt.init()
	if err != nil {
		logger.Fatal(err)
	}

	p := &ConnPool{
		opt:       opt,
		queue:     make(chan struct{}, opt.RPoolSize),
		freeConns: make([]*Conn, 0, opt.RPoolSize),
		log:       logger,
	}
	for i := 0; i < opt.RPoolSize; i++ {
		p.queue <- struct{}{}
	}

	//首先创建一个连接测试
	conn, err := p.Get()
	if err != nil {
		p.log.Fatal(err)
	}
	//归还连接测试
	p.Put(conn, false)

	//定时检测不活跃的连接
	if opt.RIdleTimeout > 0 && opt.RIdleCheckFrequency > 0 {
		go p.CheckActiveConns()
	}

	return p
}

// 获取链接，如果没有则创建一个
func (p *ConnPool) Get() (*Conn, error) {
	if p.closed {
		return nil, ErrClosed
	}
	//请求计数
	atomic.AddInt32(&p.stats.Requests, 1)

	timer := timers.Get().(*time.Timer)
	if !timer.Reset(p.opt.PoolTimeout) {
		<-timer.C
	}
	defer timers.Put(timer)

	select {
	case <-timer.C:
		//超时计数
		atomic.AddInt32(&p.stats.Timeouts, 1)
		return nil, ErrPoolTimeout
	case <-p.queue:
	}

	//channel是先进先出
	//container/list可以实现先进后出，但是效率比较低
	//这里使用slice获取最新归还的连接
	var conn *Conn
	p.lock.Lock()
	l := len(p.freeConns)
	if l > 0 {
		conn = p.freeConns[l-1]
		//清除链接
		p.freeConns = p.freeConns[:l-1]
	}
	p.lock.Unlock()

	if conn != nil {
		//命中计数
		atomic.AddInt32(&p.stats.Hits, 1)
		//连接可用 直接返回
		if conn.IsActive(p.opt.RIdleTimeout) && conn.Ping() == nil {
			atomic.AddInt32(&p.stats.UseConns, 1)
			return conn, nil
		}
		p.log.Warn(errConnActive)
		//超出最大空闲时间，或链接错误
		conn.Close()
	}

	//创建新的链接
	newcn, err := NewConn(p.opt, p.log)
	if err != nil {
		p.queue <- struct{}{}
		return nil, err
	}
	//正在使用连接计数
	atomic.AddInt32(&p.stats.UseConns, 1)
	return newcn, nil
}

//使用完后归还
func (p *ConnPool) Put(conn *Conn, forceClose bool) {
	//连接错误 强制关闭
	if forceClose || p.closed {
		conn.Close()
	} else { //归还
		p.lock.Lock()
		p.freeConns = append(p.freeConns, conn)
		p.lock.Unlock()
	}
	//减少计数
	atomic.AddInt32(&p.stats.UseConns, -1)
	if p.closed {
		return
	}
	p.queue <- struct{}{}
}

// 获取空闲的连接数
func (p *ConnPool) FreeLen() int {
	var l = 0
	p.lock.Lock()
	for _, conn := range p.freeConns {
		if conn != nil {
			l += 1
		}
	}
	p.lock.Unlock()
	return l
}

//获取连接池状态信息
func (p *ConnPool) Stats() *PoolStats {
	stats := PoolStats{}
	stats.Requests = atomic.LoadInt32(&p.stats.Requests)
	stats.Hits = atomic.LoadInt32(&p.stats.Hits)
	stats.Timeouts = atomic.LoadInt32(&p.stats.Timeouts)
	stats.UseConns = atomic.LoadInt32(&p.stats.UseConns)
	stats.FreeConns = int32(p.FreeLen())
	return &stats
}

//关闭连接池
func (p *ConnPool) Close() error {
	if p.closed {
		return ErrClosed
	}

	p.lock.Lock()
	p.closed = true
	// Close all connections.
	for _, conn := range p.freeConns {
		if conn != nil {
			conn.Close()
		}
	}
	p.freeConns = nil
	//关闭channel
	close(p.queue)
	p.lock.Unlock()

	return nil
}

//定时检测不活跃的连接
func (p *ConnPool) CheckActiveConns() {
	ticker := time.NewTicker(p.opt.RIdleCheckFrequency)
	defer ticker.Stop()

	for _ = range ticker.C {
		if p.closed {
			return
		}

		var (
			idx int
			cn  *Conn
		)

		p.lock.Lock()
		for idx, cn = range p.freeConns {
			if cn.IsActive(p.opt.RIdleTimeout) {
				break
			}
			cn.Close()
			idx += 1
		}
		if idx > 0 {
			if idx == len(p.freeConns) {
				p.freeConns = p.freeConns[:0]
			} else {
				p.freeConns = append(p.freeConns[:0], p.freeConns[idx:]...)
			}
		}
		p.lock.Unlock()

		p.log.Debug(p.opt.Addr, p.Stats())
	}

}
