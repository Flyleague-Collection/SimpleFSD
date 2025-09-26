// Package config
package config

import (
	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
)

type FsdRangeLimit struct {
	RefuseOutRange bool `json:"refuse_out_range"`
	Observer       int  `json:"observer"`
	Delivery       int  `json:"delivery"`
	Ground         int  `json:"ground"`
	Tower          int  `json:"tower"`
	Approach       int  `json:"approach"`
	Center         int  `json:"center"`
	Apron          int  `json:"apron"`
	Supervisor     int  `json:"supervisor"`
	Administrator  int  `json:"administrator"`
	FSS            int  `json:"fss"`
}

func defaultFsdRangeLimitConfig() *FsdRangeLimit {
	return &FsdRangeLimit{
		RefuseOutRange: false,
		Observer:       300,
		Delivery:       20,
		Ground:         20,
		Tower:          50,
		Approach:       150,
		Center:         600,
		Apron:          20,
		Supervisor:     300,
		Administrator:  300,
		FSS:            1500,
	}
}

func (config *FsdRangeLimit) checkValid(logger log.LoggerInterface) *ValidResult {
	if config.Observer == 0 {
		logger.Warn("Observer Range is 0, if you want to disable range limit, please set it to -1")
		config.Observer = -1
	}
	if config.Delivery == 0 {
		logger.Warn("Delivery Range is 0, if you want to disable range limit, please set it to -1")
		config.Delivery = -1
	}
	if config.Ground == 0 {
		logger.Warn("Ground Range is 0, if you want to disable range limit, please set it to -1")
		config.Ground = -1
	}
	if config.Tower == 0 {
		logger.Warn("Tower Range is 0, if you want to disable range limit, please set it to -1")
		config.Tower = -1
	}
	if config.Approach == 0 {
		logger.Warn("Approach Range is 0, if you want to disable range limit, please set it to -1")
		config.Approach = -1
	}
	if config.Center == 0 {
		logger.Warn("Center Range is 0, if you want to disable range limit, please set it to -1")
		config.Center = -1
	}
	if config.FSS == 0 {
		logger.Warn("FSS Range is 0, if you want to disable range limit, please set it to -1")
		config.FSS = -1
	}
	return ValidPass()
}
