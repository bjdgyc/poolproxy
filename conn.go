package main

import (
	"bufio"
	"net"
	"time"

	"poolproxy/slog"
)

// 读取连接的channel类型
type ChanBuf struct {
	Byte []byte
	Err  error
}

// 链接处理的接口
type Conner interface {
	Ping() error
	Auth(string, string) error
	ReadData()
	SwapData(local net.Conn) bool
}

type Conn struct {
	RawConn   *net.TCPConn
	BufReader *bufio.Reader
	Cn        Conner
	ChanRead  chan *ChanBuf
	UsedAt    time.Time
	log       *slog.Logger
}

// 创建新的远程连接
// 如果配置文件包含密码
// 需要进行密码验证
func NewConn(opt Option, logger *slog.Logger) (*Conn, error) {
	netConn, err := net.DialTimeout("tcp", opt.RAddr, 5*time.Second)
	if err != nil {
		return nil, err
	}
	tcpConn := netConn.(*net.TCPConn)

	// 设置开启KeepAlive
	tcpConn.SetKeepAlive(true)
	tcpConn.SetKeepAlivePeriod(opt.RKeepAlivePeriod)

	conn := &Conn{
		BufReader: bufio.NewReader(tcpConn),
		RawConn:   tcpConn,
		UsedAt:    time.Now(),
		ChanRead:  make(chan *ChanBuf, 1),
		log:       logger,
	}

	conn.Cn = &Redis{
		conn: conn,
	}

	// 读取数据写入channel
	go conn.Cn.ReadData()

	err = conn.Cn.Auth(opt.RUser, opt.RPass)
	if err != nil {
		conn.Close()
		return nil, err
	}

	return conn, nil
}

// 查询连接是否活跃
func (conn *Conn) IsActive(timeout time.Duration) bool {
	return timeout > 0 && time.Since(conn.UsedAt) < timeout
}

// 获取channel chan
func (conn *Conn) GetReadChan() <-chan *ChanBuf {
	return conn.ChanRead
}

func (conn *Conn) Write(b []byte) error {
	_, err := conn.RawConn.Write(b)
	conn.UsedAt = time.Now()
	return err
}

// 远程地址
func (conn *Conn) RemoteAddr() net.Addr {
	return conn.RawConn.RemoteAddr()
}

func (conn *Conn) Close() error {
	err := conn.RawConn.Close()
	return err
}

// redis的ping
func (conn *Conn) Ping() error {
	return conn.Cn.Ping()
}

// 交换数据
func (conn *Conn) SwapData(local net.Conn) bool {
	return conn.Cn.SwapData(local)
}
