package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

type Redis struct{
	conn *Conn
}

var _ Conner = (*Redis)(nil)

//redis的ping
func (cn *Redis) Ping() error {
	err := cn.conn.Write([]byte("*1\r\n$4\r\nPING\r\n"))
	if err != nil {
		return nil
	}
	cb := <-cn.conn.GetReadChan()
	if cb.Err != nil {
		return cb.Err
	}
	if strings.ToUpper(string(cb.Byte)) != "+PONG\r\n" {
		return fmt.Errorf("error ping")
	}
	return nil
}

//redis 权限验证
func (cn *Redis) Auth(user, pass string) error {
	if pass == "" {
		return nil
	}

	data := fmt.Sprintf("*2\r\n$4\r\nAUTH\r\n$%d\r\n%s\r\n", len(pass), pass)
	err := cn.conn.Write([]byte(data))
	if err != nil {
		return err
	}
	cb := <-cn.conn.GetReadChan()
	if cb.Err != nil {
		return cb.Err
	}
	if strings.ToUpper(string(cb.Byte)) != "+OK\r\n" {
		return fmt.Errorf("auth error")
	}
	return nil
}

//redis 数据读取
func (cn *Redis) ReadData() {
	var (
		line []byte
		err  error
		cb   *ChanBuf
	)
	for {
		line, err = cn.conn.BufReader.ReadBytes('\n')
		cn.conn.UsedAt = time.Now()
		cb = &ChanBuf{Byte: line, Err: err}
		cn.conn.ChanRead <- cb
		if err != nil {
			break
		}
	}
}

//redis数据交换
func (cn *Redis) SwapData(local net.Conn) bool {
	lread := bufio.NewReader(local)

	var (
		err           error
		forceClose    = false
		exitChanProxy = make(chan struct{})
		cb            *ChanBuf
	)

	//读取数据
	go func() {
		var (
			line []byte
			err  error
		)
		for {
			line, err = lread.ReadBytes('\n')
			if err != nil {
				if err != io.EOF {
					cn.conn.log.Error("local read error:", err)
				}
				break
			}
			fmt.Println(string(line))
			err = cn.conn.Write(line)
			if err != nil {
				forceClose = true
				cn.conn.log.Error("remote write error:", err)
				break
			}
		}
		exitChanProxy <- struct{}{}
	}()

	//客户端写回数据
	readChan := cn.conn.GetReadChan()
	for {
		select {
		case <-exitChanProxy:
			goto FAIL
		case cb = <-readChan:
			if cb.Err != nil {
				forceClose = true
				cn.conn.log.Error("remote read error:", err)
				goto FAIL
			}
			//fmt.Println(string(cb.Byte))
			_, err = local.Write(cb.Byte)
			if err != nil {
				cn.conn.log.Error("local write error:", err)
				goto FAIL
			}
		}
	}
FAIL:
	return forceClose
}
