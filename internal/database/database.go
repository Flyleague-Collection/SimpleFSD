package database

import (
	"bufio"
	"context"
	"github.com/fsnotify/fsnotify"
	"github.com/half-nothing/simple-fsd/internal/interfaces/config"
	"github.com/half-nothing/simple-fsd/internal/interfaces/global"
	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
	. "github.com/half-nothing/simple-fsd/internal/interfaces/operation"
	"github.com/half-nothing/simple-fsd/internal/utils"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type DBCloseCallback struct {
	logger  log.LoggerInterface
	watcher *fsnotify.Watcher
}

func NewDBCloseCallback(logger log.LoggerInterface, watcher *fsnotify.Watcher) *DBCloseCallback {
	return &DBCloseCallback{
		logger:  logger,
		watcher: watcher,
	}
}

func (dc *DBCloseCallback) Invoke(_ context.Context) error {
	dc.logger.InfoF("Closing file watcher")
	return dc.watcher.Close()
}

var (
	data = map[string]*User{}
	lock = sync.RWMutex{}
)

func readData(file *os.File) {
	lock.Lock()
	defer lock.Unlock()

	data = map[string]*User{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := string(scanner.Bytes())
		if strings.HasPrefix(line, "#") {
			continue
		}
		information := strings.Split(line, " ")
		if len(information) < 3 {
			continue
		}
		user := &User{
			Cid:      information[0],
			Password: information[1],
			Rating:   utils.StrToInt(information[2], -1),
		}
		data[user.Cid] = user
	}
}

func ConnectDatabase(lg log.LoggerInterface, config *config.Config, _ bool) (*DBCloseCallback, *DatabaseOperations, error) {
	if err := os.MkdirAll(filepath.Dir(config.CertFile), global.DefaultDirectoryPermission); err != nil {
		return nil, nil, err
	}

	var file *os.File
	if _, err := os.Stat(config.CertFile); os.IsNotExist(err) {
		file, _ = os.OpenFile(config.CertFile, os.O_WRONLY|os.O_CREATE, global.DefaultFilePermissions)
		data := "######################\n" +
			"# -1 Ban\n" +
			"# 0 Normal\n" +
			"# 1 Observer\n" +
			"# 2 STU1\n" +
			"# 3 STU2\n" +
			"# 4 STU3\n" +
			"# 5 CTR1\n" +
			"# 6 CTR2\n" +
			"# 7 CTR3\n" +
			"# 8 Instructor1\n" +
			"# 9 Instructor2\n" +
			"# 10 Instructor3\n" +
			"# 11 Supervisor\n" +
			"# 12 Administrator\n" +
			"######################\n" +
			"# CID PASSWORD RATING"
		_, _ = file.Write([]byte(data))
	} else if err != nil {
		return nil, nil, err
	} else if file, err = os.Open(config.CertFile); err != nil {
		return nil, nil, err
	}

	readData(file)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, nil, err
	}
	if err := watcher.Add(config.CertFile); err != nil {
		return nil, nil, err
	}

	go func() {
		for {
			select {
			case ev := <-watcher.Events:
				{
					if ev.Op&fsnotify.Write == fsnotify.Write {
						file, err := os.Open(ev.Name)
						if err != nil {
							lg.ErrorF("Error opening file %s", ev.Name)
						} else {
							readData(file)
						}
					}
				}
			case err := <-watcher.Errors:
				{
					lg.ErrorF("Error watching file, %v", err)
					return
				}
			}
		}
	}()

	return NewDBCloseCallback(lg, watcher), NewDatabaseOperations(NewUserOperation(lg, config), NewFlightPlanOperation(lg, config)), nil
}
