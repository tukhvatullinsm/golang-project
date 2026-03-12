package server

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"
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

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
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
		wa.Objlog.Info("Error writing response:", err)
	}

}

func (wa *WebApp) GetValueJSON(w http.ResponseWriter, r *http.Request) {
	objMetrics := new(Metrics)
	w.Header().Set("Content-Type", "application/json")
	switch r.Header.Get("Content-Type") {
	case "application/json":
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&objMetrics)
		if err != nil {
			wa.Objlog.Info("Error decoding JSON:", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		switch objMetrics.MType {
		case "counter":
			objMetrics.Delta = new(int64)
			res := wa.ObjStorage.GetValue(objMetrics.MType, objMetrics.ID)
			if res == nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			*objMetrics.Delta, _ = strconv.ParseInt(fmt.Sprintf("%v", res), 10, 64)
		case "gauge":
			objMetrics.Value = new(float64)
			res := wa.ObjStorage.GetValue(objMetrics.MType, objMetrics.ID)
			if res == nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusOK)
			*objMetrics.Value, _ = strconv.ParseFloat(fmt.Sprintf("%v", res), 64)
		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.NewEncoder(w).Encode(objMetrics)
		if err != nil {
			wa.Objlog.Info("Error encoding JSON:", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
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

func (wa *WebApp) SetValuesJSON(w http.ResponseWriter, r *http.Request) {
	var objMetrics Metrics
	w.Header().Set("Content-Type", "application/json")
	switch r.Header.Get("Content-Type") {
	case "application/json":
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&objMetrics)
		if err != nil {
			wa.Objlog.Infoln("Error decoding JSON:", err,
				r.Body)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		switch objMetrics.MType {
		case "counter":
			wa.ObjStorage.SetValue(objMetrics.MType, objMetrics.ID, strconv.FormatInt(*objMetrics.Delta, 10))
			res := wa.ObjStorage.GetValue(objMetrics.MType, objMetrics.ID)
			*objMetrics.Delta, _ = strconv.ParseInt(fmt.Sprintf("%v", res), 10, 64)
		case "gauge":
			wa.ObjStorage.SetValue(objMetrics.MType, objMetrics.ID, strconv.FormatFloat(*objMetrics.Value, 'f', -1, 64))
			res := wa.ObjStorage.GetValue(objMetrics.MType, objMetrics.ID)
			*objMetrics.Value, _ = strconv.ParseFloat(fmt.Sprintf("%v", res), 64)
		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.NewEncoder(w).Encode(objMetrics)
		if err != nil {
			wa.Objlog.Info("Error encoding JSON:", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if ok := slices.Contains(wa.Parameters, objMetrics.ID); !ok {
		wa.Parameters = append(wa.Parameters, objMetrics.ID)
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

func (wa *WebApp) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
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

func (wa *WebApp) GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") == "html/text" || r.Header.Get("Content-Type") == "application/json" {
			if r.Header.Get("Content-Encoding") == "gzip" {
				gz, err := gzip.NewReader(r.Body)
				if err != nil {
					wa.Objlog.Info("Error creating gzip reader:", err)
					http.Error(w, "Internal server error", http.StatusInternalServerError)
					return
				}
				defer gz.Close()
				r.Body.Close()
				uncomBody, err := io.ReadAll(gz)
				if err != nil {
					wa.Objlog.Info("Error creating gzip reader:", err)
				}
				newBody := io.NopCloser(bytes.NewBuffer(uncomBody))
				r.Body = newBody
				r.ContentLength = int64(len(uncomBody))
			}
			if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
				if err != nil {
					wa.Objlog.Info("Error creating gzip writer:", err)
					io.WriteString(w, err.Error())
					return
				}
				defer gz.Close()
				w.Header().Add("Content-Encoding", "gzip")
				w = gzipWriter{ResponseWriter: w, Writer: gz}
			}
		}
		next.ServeHTTP(w, r)

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

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
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
