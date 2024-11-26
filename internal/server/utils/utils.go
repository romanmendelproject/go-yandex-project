package utils

import (
	"errors"
	"strings"
)

type UrlParamType struct {
	MapName string
	Limit   string
}

func ParseURLMapValue(url string) (UrlParamType, error) {
	var urlParam UrlParamType

	urlData := strings.Split(url[1:], "/")
	if len(urlData) < 4 {
		return UrlParamType{}, errors.New("error parameters from URL")
	}

	urlParam.MapName = urlData[3]

	if len(urlData) == 5 {
		limitReg := strings.Split(urlData[4], "limit=")
		if len(limitReg) != 2 {
			return UrlParamType{}, errors.New("error parameters limit from URL")
		}
		urlParam.Limit = limitReg[0]
	}

	return urlParam, nil
}

func GetFloatPtr(v float64) *float64 {
	return &v
}
