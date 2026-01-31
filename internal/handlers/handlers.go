package handlers

import (
	"net/http"
	"strconv"
)

var parameters = map[string]int64{
	"counter": 1,
	"gauge":   2,
}

func CheckUrlParams(params []string) (bool, int) {
	if len(params) != 3 {
		return false, http.StatusNotFound
	}
	if _, ok := parameters[params[0]]; !ok {
		return false, http.StatusBadRequest
	}
	if _, err := strconv.ParseInt(params[2], 10, 64); err != nil {
		return false, http.StatusBadRequest
	}
	return true, 0
}
