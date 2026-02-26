package json

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

var (
	MissingBody = errors.New("Missing Body")
)

func ParseJSON(r *http.Request, payload any) error {
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		return errors.New("Invalid Content-Type header")
	}

	if r.Body == nil || r.ContentLength == 0 {
		return MissingBody
	}

	return json.NewDecoder(r.Body).Decode(payload)
}

func SendJSON(w http.ResponseWriter, status int, body any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(body)
}
