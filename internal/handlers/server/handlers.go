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
	GetAllValue() *map[string]any
}

var parameters = map[string]int64{
	"counter": 1,
	"gauge":   2,
}

type WebApp struct {
	ObjStorage IntMemStorage
	Parameters []string
}

func (wa *WebApp) Init(stg *storage.MemStorage) {
	wa.ObjStorage = stg
	wa.Parameters = make([]string, 0)
}

/*
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
				wa.ObjStorage.SetValue(sliceQueryParams[0], sliceQueryParams[1], sliceQueryParams[2])
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
*/
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
	//ct := r.Header.Get("Content-Type")
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

	/*switch ct {
	case "text/plain":
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
		return
	default:
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	} */
}

func (wa *WebApp) GetParam(w http.ResponseWriter, r *http.Request) {
	listObj := *wa.ObjStorage.GetAllValue()
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
