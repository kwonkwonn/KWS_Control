package util

import (
	"encoding/json"
	"io"
	"net/http"
)

func UnmarshalBodyAndClose[T any](body io.ReadCloser) (T, error) {
	var t T
	err := json.NewDecoder(body).Decode(&t)

	e := body.Close()
	if err == nil {
		err = e
	}

	return t, err
}

// 메서드 확인 후, 일치하지 않으면 http.Error 반환
func CheckMethod(w http.ResponseWriter, r *http.Request, expectedMethod string) bool {
	if r.Method != expectedMethod {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return false
	}
	return true
}

