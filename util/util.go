package util

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

func CheckMethod(w http.ResponseWriter, r *http.Request, expectedMethod string) bool {
	if r.Method != expectedMethod {
		log := logrus.New()
		log.SetReportCaller(true)
		
		h := w.Header()
		h.Del("Content-Length")
		h.Set("Content-Type", "text/plain; charset=utf-8")
		h.Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(http.StatusMethodNotAllowed)

		log.Errorf("Invalid request method: %s", r.Method)

		return false
	}
	return true
}

