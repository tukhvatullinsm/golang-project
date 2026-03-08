package server

import (
	"fmt"
	"log"
	"net/http"
	"slices"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/tukhvatullinsm/golang-project/internal/storage"
)

type IntMemStorage interface {
	SetValue(param, key, value string)
	GetValue(param, key string) any
	GetAllValue() map[string]any
}

var parameters = map[string]struct{}{
	"counter": {},
	"gauge":   {},
}

type WebApp struct {
	ObjStorage IntMemStorage
	Parameters []string
}

func (wa *WebApp) Init(stg *storage.MemStorage) {
	wa.ObjStorage = stg
	wa.Parameters = make([]string, 0)
}

func (wa *WebApp) GetValue(w http.ResponseWriter, r *http.Request) {
	typeAtt := chi.URLParam(r, "type")
	nameAtt := chi.URLParam(r, "name")
	w.Header().Set("Content-Type", "text/plain")
	res := wa.ObjStorage.GetValue(typeAtt, nameAtt)
	if res == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(fmt.Sprintf("%v", res)))
	if err != nil {
		log.Fatal("Error writing response:", err)
	}

}

func (wa *WebApp) SetValues(w http.ResponseWriter, r *http.Request) {
	typeAtt := chi.URLParam(r, "type")
	nameAtt := chi.URLParam(r, "name")
	valueAtt := chi.URLParam(r, "value")
	if typeAtt == "" || nameAtt == "" || valueAtt == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	res, httpCode := CheckURLParams(typeAtt, valueAtt)
	if !res {
		w.WriteHeader(httpCode)
		return
	}
	wa.ObjStorage.SetValue(typeAtt, nameAtt, valueAtt)
	if ok := slices.Contains(wa.Parameters, nameAtt); !ok {
		wa.Parameters = append(wa.Parameters, nameAtt)
	}
	w.WriteHeader(http.StatusOK)
}

func (wa *WebApp) GetParam(w http.ResponseWriter, r *http.Request) {
	listObj := wa.ObjStorage.GetAllValue()
	rawResult := ""
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html")
	rawResult += "<p>Метрики системы</p>\n"
	rawResult += "<ul>\n"
	for _, v := range wa.Parameters {
		rawResult += fmt.Sprintf("<li>%s : %s</li>\n", v, fmt.Sprintf("%v", listObj[v]))
	}
	rawResult += "</ul>\n"
	_, err := w.Write([]byte(rawResult))
	if err != nil {
		log.Fatal("Error writing response:", err)
	}
}

func CheckURLParams(params ...string) (bool, int) {
	if _, ok := parameters[params[0]]; !ok {
		return false, http.StatusBadRequest
	}
	switch params[0] {
	case "gauge":
		if _, err := strconv.ParseFloat(params[1], 64); err != nil {
			return false, http.StatusBadRequest
		}
	case "counter":
		if _, err := strconv.ParseInt(params[1], 10, 64); err != nil {
			return false, http.StatusBadRequest
		}
	}
	return true, 0
}
