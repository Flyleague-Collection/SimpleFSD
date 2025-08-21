package fsd_server

import (
	"context"
	c "github.com/half-nothing/fsd-server/internal/config"
	"github.com/half-nothing/fsd-server/internal/fsd_server/packet"
	"github.com/half-nothing/fsd-server/internal/interfaces/fsd"
	"net"
	"time"
)

type FsdCloseCallback struct {
	clientManager fsd.ClientManagerInterface
}

func NewFsdCloseCallback(clientManager fsd.ClientManagerInterface) *FsdCloseCallback {
	return &FsdCloseCallback{clientManager: clientManager}
}

func (dc *FsdCloseCallback) Invoke(ctx context.Context) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	done := make(chan struct{})
	go func() {
		if err := dc.clientManager.Shutdown(timeoutCtx); err != nil {
			return
		}
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-timeoutCtx.Done():
		return timeoutCtx.Err()
	}
}

// StartFSDServer 启动FSD服务器
func StartFSDServer(config *c.Config) {
	// 初始化客户端管理器
	cm := packet.NewClientManager(config)

	// 创建TCP监听器
	sem := make(chan struct{}, config.Server.FSDServer.MaxWorkers)
	ln, err := net.Listen("tcp", config.Server.FSDServer.Address)
	if err != nil {
		c.FatalF("FSD Server Start error: %v", err)
		return
	}
	c.InfoF("FSD Server Listen On " + ln.Addr().String())

	// 确保在函数退出时关闭监听器
	defer func() {
		err := ln.Close()
		if err != nil {
			c.ErrorF("Server close error: %v", err)
		}
	}()

	c.GetCleaner().Add(NewFsdCloseCallback(cm))

	// 循环接受新的连接
	for {
		conn, err := ln.Accept()
		if err != nil {
			c.ErrorF("Accept connection error: %v", err)
			continue
		}

		c.DebugF("Accepted new connection from %s", conn.RemoteAddr().String())

		// 使用信号量控制并发连接数
		sem <- struct{}{}
		go func(c net.Conn) {
			connection := packet.NewConnectionHandler(conn, conn.RemoteAddr().String(), config, cm)
			connection.HandleConnection()
			// 释放信号量
			<-sem
		}(conn)
	}
}
