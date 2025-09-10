package base

import (
	"context"
	"fmt"
	. "github.com/half-nothing/simple-fsd/internal/interfaces/global"
	. "github.com/half-nothing/simple-fsd/internal/interfaces/log"
	"github.com/half-nothing/simple-fsd/internal/utils"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Cleaner struct {
	cleaners       []Callable
	mu             sync.Mutex
	cleaning       bool
	loggerShutdown Callable
	logger         LoggerInterface
}

func NewCleaner(logger LoggerInterface) *Cleaner {
	return &Cleaner{
		cleaners:       make([]Callable, 0),
		loggerShutdown: logger.ShutdownCallback(),
		logger:         logger,
	}
}

func (c *Cleaner) Add(callable Callable) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.cleaning {
		c.logger.Debug("Cleaner is already shutting down, ignoring new cleaner")
		return
	}
	c.cleaners = append(c.cleaners, callable)
	c.logger.DebugF("Adding cleaner #%d (%T)", len(c.cleaners), callable)
}

func (c *Cleaner) Clean() {
	c.mu.Lock()
	c.cleaning = true // 标记为清理中，阻止后续Add操作
	cleanersCopy := make([]Callable, len(c.cleaners))
	copy(cleanersCopy, c.cleaners)
	c.mu.Unlock()

	c.logger.DebugF("Starting cleanup of %d registered functions", len(cleanersCopy))

	var errs []error
	utils.ReverseForEach(cleanersCopy, func(idx int, callback Callable) { // 使用匿名函数确保defer在每次迭代执行
		c.logger.DebugF("Invoking cleaner #%d (%T)", idx+1, callback)
		timeoutCtx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelFunc() // 确保每次调用后取消上下文
		if err := callback.Invoke(timeoutCtx); err != nil {
			c.logger.ErrorF("Cleaner #%d (%T) failed: %v", idx+1, callback, err) // 记录类型和错误
			errs = append(errs, err)
		}
	})

	if len(errs) > 0 {
		c.logger.ErrorF("%d errors occurred during cleanup:", len(errs))
		for i, err := range errs {
			c.logger.ErrorF("Error %d: %v", i+1, err)
		}
	} else {
		c.logger.Debug("All cleaners executed successfully")
	}
	c.logger.Info("Cleanup finished, server offline")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := c.loggerShutdown.Invoke(shutdownCtx); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "LOGGER SHUTDOWN ERROR: %v\n", err)
	}
	syscall.Exit(0)
}

func (c *Cleaner) Init() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ctx.Done()
		stop()
		c.logger.Info("Received interrupt signal, shutting down")

		c.Clean()
	}()
}
