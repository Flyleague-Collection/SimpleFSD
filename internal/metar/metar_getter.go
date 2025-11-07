// Package metar
package metar

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	. "github.com/half-nothing/simple-fsd/internal/interfaces"
	"github.com/half-nothing/simple-fsd/internal/interfaces/config"
	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
	"github.com/mdaverde/jsonpath"
)

type MetarGetter struct {
	logger log.LoggerInterface
	config *config.MetarSource
}

func NewMetarGetter(
	logger log.LoggerInterface,
	config *config.MetarSource,
) *MetarGetter {
	return &MetarGetter{
		logger: logger,
		config: config,
	}
}

func (getter *MetarGetter) GetMetar(icao string) (string, error) {
	url := fmt.Sprintf(getter.config.Url, icao)

	getter.logger.InfoF("Get metar from url %s", url)
	response, err := http.Get(url)
	if err != nil {
		getter.logger.ErrorF("Get metar from url %s failed: %s", url, err.Error())
		return "", err
	}

	if response.StatusCode != 200 {
		getter.logger.ErrorF("Get metar from url %s failed: status code %d", url, response.StatusCode)
		return "", errors.New("status code is not 200")
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		getter.logger.ErrorF("Get metar from url %s failed: %s", url, err.Error())
		return "", err
	}
	_ = response.Body.Close()
	data = bytes.TrimRight(data, "\n")

	switch getter.config.MetarSourceType {
	case config.Raw:
		if getter.config.Multiline == "" {
			return string(data), nil
		}
		metars := strings.Split(string(data), getter.config.Multiline)
		if getter.config.Reverse {
			return metars[len(metars)-1], nil
		}
		return metars[0], nil
	case config.Json:
		var jsonData interface{}
		err := json.Unmarshal(data, &jsonData)
		if err != nil {
			getter.logger.ErrorF("Get metar from url %s failed: %s", url, err.Error())
			return "", err
		}
		value, err := jsonpath.Get(&jsonData, getter.config.Selector)
		if err != nil {
			getter.logger.ErrorF("Get metar from url %s failed: %s", url, err.Error())
			return "", err
		}
		if val, ok := value.(string); ok {
			return val, nil
		}
		getter.logger.ErrorF("Get metar from url %s failed: value is not string", url)
		return "", errors.New("target is not a string")
	case config.Html:
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
		if err != nil {
			getter.logger.ErrorF("Get metar from url %s failed: %s", url, err.Error())
			return "", err
		}
		rawData := doc.Find(getter.config.Selector)
		if rawData.Size() == 0 {
			getter.logger.ErrorF("Get metar from url %s failed: no metar found", url)
			return "", ErrMetarNotFound
		}
		if getter.config.Multiline == "" {
			if getter.config.Reverse {
				return rawData.Get(rawData.Size() - 1).FirstChild.Data, nil
			}
			return rawData.Get(0).FirstChild.Data, nil
		}
		firstChild := rawData.Get(0).FirstChild
		if firstChild == nil {
			getter.logger.ErrorF("Get metar from url %s failed: no first child node found", url)
			return "", ErrMetarNotFound
		}
		metars := strings.Split(firstChild.Data, getter.config.Multiline)
		if len(metars) < 1 {
			getter.logger.ErrorF("Get metar from url %s failed: no metar found", url)
			return "", ErrMetarNotFound
		}
		if getter.config.Reverse {
			return metars[len(metars)-1], nil
		}
		return metars[0], nil
	}
	return "", ErrMetarNotFound
}
