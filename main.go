package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/eddwinpaz/checkout-logging/logging"
	mylog "github.com/sirupsen/logrus"
)

type Response struct {
	Description string `json:"description"`
	Status      bool   `json:"status"`
}

func sampleHandler(w http.ResponseWriter, r *http.Request) {

	var response = Response{
		Description: "hello again",
		Status:      true,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)

}

type logResponseWriter struct {
	status int
	body   string
	http.ResponseWriter
}

func (w *logResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *logResponseWriter) Write(body []byte) (int, error) {
	w.body = string(body)
	return w.ResponseWriter.Write(body)
}

func responseLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loggingRW := &logResponseWriter{
			ResponseWriter: w,
		}
		logging.InitializeLogging("vendors.log")
		startTime := time.Now()
		defer func() {
			duration := time.Since(startTime)
			body, _ := ioutil.ReadAll(r.Body)
			mylog.Infof("body=%s method=%s path=%s duration=%f response=%s", body, r.Method, r.URL.Path, duration.Seconds(), loggingRW.body)
		}()
		h.ServeHTTP(loggingRW, r)
	})
}

func main() {
	http.Handle("/", responseLogger(http.HandlerFunc(sampleHandler)))
	http.ListenAndServe(":9000", nil)
}
