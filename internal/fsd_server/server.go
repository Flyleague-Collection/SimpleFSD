package fsd_server

import (
	"context"
	"github.com/half-nothing/simple-fsd/internal/fsd_server/packet"
	. "github.com/half-nothing/simple-fsd/internal/interfaces"
	"github.com/half-nothing/simple-fsd/internal/interfaces/fsd"
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
func StartFSDServer(applicationContent *ApplicationContent) {
	config := applicationContent.ConfigManager().Config()
	logger := applicationContent.Logger().FsdLogger()

	// 创建TCP监听器
	sem := make(chan struct{}, config.Server.FSDServer.MaxWorkers)
	ln, err := net.Listen("tcp", config.Server.FSDServer.Address)
	if err != nil {
		logger.FatalF("FSD Server Start error: %v", err)
		return
	}
	logger.InfoF("FSD Server Listen On " + ln.Addr().String())

	// 确保在函数退出时关闭监听器
	defer func() {
		err := ln.Close()
		if err != nil {
			logger.ErrorF("Server close error: %v", err)
		}
	}()

	applicationContent.Cleaner().Add(NewFsdCloseCallback(applicationContent.ClientManager()))

	// 循环接受新的连接
	for {
		conn, err := ln.Accept()
		if err != nil {
			logger.ErrorF("Accept connection error: %v", err)
			continue
		}

		logger.DebugF("Accepted new connection from %s", conn.RemoteAddr().String())

		// 使用信号量控制并发连接数
		sem <- struct{}{}
		go func(c net.Conn) {
			connection := packet.NewSession(applicationContent, conn)
			connection.HandleConnection(config.Server.FSDServer.HeartbeatDuration)
			// 释放信号量
			<-sem
		}(conn)
	}
}
