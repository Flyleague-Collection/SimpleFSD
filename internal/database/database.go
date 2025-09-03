package database

import (
	"bufio"
	"context"
	"github.com/fsnotify/fsnotify"
	c "github.com/half-nothing/simple-fsd/internal/config"
	. "github.com/half-nothing/simple-fsd/internal/interfaces/operation"
	"github.com/half-nothing/simple-fsd/internal/utils"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type DBCloseCallback struct {
	watcher *fsnotify.Watcher
}

func NewDBCloseCallback(watcher *fsnotify.Watcher) *DBCloseCallback {
	return &DBCloseCallback{watcher: watcher}
}

func (dc *DBCloseCallback) Invoke(_ context.Context) error {
	c.InfoF("Closing file watcher")
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

func ConnectDatabase(config *c.Config) (*DatabaseOperations, error) {
	if err := os.MkdirAll(filepath.Dir(config.CertFile), 0775); err != nil {
		return nil, err
	}

	var file *os.File
	if _, err := os.Stat(config.CertFile); os.IsNotExist(err) {
		file, _ = os.Create(config.CertFile)
		_, _ = file.Write([]byte("# CID PASSWORD RATING"))
	} else if err != nil {
		return nil, err
	} else if file, err = os.Open(config.CertFile); err != nil {
		return nil, err
	}

	readData(file)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	if err := watcher.Add(config.CertFile); err != nil {
		return nil, err
	}

	go func() {
		for {
			select {
			case ev := <-watcher.Events:
				{
					if ev.Op&fsnotify.Write == fsnotify.Write {
						file, err := os.Open(ev.Name)
						if err != nil {
							c.ErrorF("Error opening file %s", ev.Name)
						} else {
							readData(file)
						}
					}
				}
			case err := <-watcher.Errors:
				{
					c.ErrorF("Error watching file, %v", err)
					return
				}
			}
		}
	}()

	c.GetCleaner().Add(NewDBCloseCallback(watcher))

	return NewDatabaseOperations(NewUserOperation(config), NewFlightPlanOperation(config)), nil
}
