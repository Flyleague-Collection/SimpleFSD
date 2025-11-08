// Package config
package config

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/half-nothing/simple-fsd/src/interfaces/logger"
)

type MetarSourceType string

const (
	Raw  MetarSourceType = "raw"
	Json MetarSourceType = "json"
	Html MetarSourceType = "html"
)

var allowedMetarSources = []MetarSourceType{Raw, Json, Html}

type MetarSource struct {
	Url             string          `json:"url"`
	ReturnType      string          `json:"return_type"`
	MetarSourceType MetarSourceType `json:"-"`
	Reverse         bool            `json:"reverse"`
	Selector        string          `json:"selector"`
	Multiline       string          `json:"multiline"`
}

type MetarSources []*MetarSource

func defaultMetarSources() MetarSources {
	return []*MetarSource{
		{
			Url:        "https://aviationweather.gov/api/data/metar?ids=%s",
			ReturnType: "raw",
			Reverse:    false,
			Selector:   "",
			Multiline:  "",
		},
	}
}

func (config MetarSources) checkValid(logger logger.LoggerInterface) *ValidResult {
	for _, metarSource := range config {
		metarSource.MetarSourceType = MetarSourceType(strings.ToLower(metarSource.ReturnType))
		if !slices.Contains(allowedMetarSources, metarSource.MetarSourceType) {
			return ValidFail(fmt.Errorf("unsupported metar source type: %s", metarSource.ReturnType))
		}
		switch metarSource.MetarSourceType {
		case Raw:
			if metarSource.Reverse && metarSource.Multiline == "" {
				logger.Warn("when set multiline to empty, reverse dont take effect")
			}
		case Json:
			if metarSource.Selector == "" {
				return ValidFail(errors.New("when set return_type to json, selector cannot be empty"))
			}
			if metarSource.Reverse {
				logger.Warn("when set return_type to json, reverse dont take effect")
			}
			if metarSource.Multiline != "" {
				logger.Warn("when set return_type to json, multiline dont take effect")
			}
		case Html:
			if metarSource.Selector == "" {
				return ValidFail(errors.New("when set return_type to html, selector cannot be empty"))
			}
			if metarSource.Reverse && metarSource.Multiline == "" {
				logger.Warn("when set multiline to empty, reverse dont take effect")
			}
		}
	}
	return ValidPass()
}
