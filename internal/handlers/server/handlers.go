package server

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/tukhvatullinsm/golang-project/internal/storage"
)

type WebApp struct {
	ObjStorage *storage.MemStorage
}

func (wa *WebApp) Init(strg *storage.MemStorage) {
	wa.ObjStorage = strg
}

func (wa *WebApp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		rawQueryParams := r.URL.EscapedPath()
		strippedQueryParams := strings.TrimPrefix(rawQueryParams, "/update/")
		sliceQueryParams := strings.Split(strippedQueryParams, "/")
		ct := r.Header.Get("Content-Type")
		switch ct {
		case "text/plain":
			res, httpCode := CheckURLParams(sliceQueryParams)
			if !res {
				w.WriteHeader(httpCode)
				return
			}
			wa.ObjStorage.Set(sliceQueryParams[0], sliceQueryParams[1], sliceQueryParams[2])
			w.WriteHeader(http.StatusOK)
			return
		default:
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}
	}
	if r.Method == "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

var parameters = map[string]int64{
	"counter": 1,
	"gauge":   2,
}

func CheckURLParams(params []string) (bool, int) {
	if len(params) != 3 {
		return false, http.StatusNotFound
	}
	if _, ok := parameters[params[0]]; !ok {
		return false, http.StatusBadRequest
	}
	switch params[0] {
	case "gauge":
		if _, err := strconv.ParseFloat(params[2], 64); err != nil {
			return false, http.StatusBadRequest
		}
	case "counter":
		if _, err := strconv.ParseInt(params[2], 10, 64); err != nil {
			return false, http.StatusBadRequest
		}
	}
	return true, 0
}
