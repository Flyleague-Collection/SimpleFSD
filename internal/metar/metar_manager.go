// Package metar
package metar

import (
	"sync"
	"time"

	. "github.com/half-nothing/simple-fsd/internal/interfaces"
	"github.com/half-nothing/simple-fsd/internal/interfaces/config"
	"github.com/half-nothing/simple-fsd/internal/interfaces/global"
	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
	"golang.org/x/sync/singleflight"
)

type MetarManager struct {
	logger       log.LoggerInterface
	config       config.MetarSources
	getters      []MetarGetterInterface
	metarCache   CacheInterface[*string]
	requestGroup singleflight.Group
}

func NewMetarManager(
	logger log.LoggerInterface,
	config config.MetarSources,
	cache CacheInterface[*string],
) *MetarManager {
	manager := &MetarManager{
		logger:     log.NewLoggerAdapter(logger, "MetarManager"),
		config:     config,
		getters:    make([]MetarGetterInterface, 0, len(config)),
		metarCache: cache,
	}
	for _, metarSource := range config {
		manager.getters = append(manager.getters, NewMetarGetter(manager.logger, metarSource))
	}
	return manager
}

func (metarManager *MetarManager) cacheMetar(icao string, metar *string) {
	currentTime := time.Now()
	minute := currentTime.Minute()
	var addMinutes int
	if minute < 30 {
		addMinutes = 30 - minute
	} else {
		addMinutes = 60 - minute
	}
	next := currentTime.Add(time.Duration(addMinutes) * time.Minute)
	expirationTime := time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), next.Minute(), 0, 0, next.Location())
	metarManager.metarCache.Set(icao, metar, expirationTime)
}

func (metarManager *MetarManager) QueryMetar(icao string) (metar string, err error) {
	if icao == "" || len(icao) != 4 {
		return "", ErrICAOInvalid
	}

	if cachedMetar, ok := metarManager.metarCache.Get(icao); ok {
		if cachedMetar == nil {
			return "", ErrMetarNotFound
		}
		return *cachedMetar, nil
	}

	result, err, _ := metarManager.requestGroup.Do(icao, func() (interface{}, error) {
		for _, getter := range metarManager.getters {
			metar, err := getter.GetMetar(icao)
			if err != nil {
				continue
			}
			metarManager.cacheMetar(icao, &metar)
			return metar, nil
		}
		metarManager.cacheMetar(icao, nil)
		return "", ErrMetarNotFound
	})
	metar = result.(string)
	return
}

func (metarManager *MetarManager) QueryMetars(icaos []string) (metars []string) {
	wg := sync.WaitGroup{}
	lock := sync.Mutex{}
	metars = make([]string, 0, len(icaos))
	limiter := make(chan struct{}, *global.MetarQueryThread)

	for _, icao := range icaos {
		wg.Add(1)
		limiter <- struct{}{}
		go func() {
			defer func() {
				<-limiter
				wg.Done()
			}()
			metar, err := metarManager.QueryMetar(icao)
			if err != nil {
				metarManager.logger.ErrorF("QueryMetars err %s", err)
				return
			}
			lock.Lock()
			metars = append(metars, metar)
			lock.Unlock()
		}()
	}
	wg.Wait()
	return
}
