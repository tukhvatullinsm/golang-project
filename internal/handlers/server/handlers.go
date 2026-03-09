package server

import (
	"fmt"
	"log"
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/tukhvatullinsm/golang-project/internal/storage"
	"go.uber.org/zap"
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
	Objlog     *zap.SugaredLogger
}

func (wa *WebApp) Init(stg *storage.MemStorage, lg *zap.SugaredLogger) {
	wa.ObjStorage = stg
	wa.Parameters = make([]string, 0)
	wa.Objlog = lg
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

func (wa *WebApp) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w, // встраиваем оригинальный http.ResponseWriter
			responseData:   responseData,
		}

		next.ServeHTTP(&lw, r)

		duration := time.Since(start)

		wa.Objlog.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"status", responseData.status,
			"duration", duration,
			"size", responseData.size,
		)
	})
}

type responseData struct {
	status int
	size   int
}

type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := lrw.ResponseWriter.Write(b)
	lrw.responseData.size += size
	return size, err
}

func (lrw *loggingResponseWriter) WriteHeader(status int) {
	lrw.responseData.status = status
	lrw.ResponseWriter.WriteHeader(status)
}
