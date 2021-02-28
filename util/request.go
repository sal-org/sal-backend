package util

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// RequiredFiledsCheck - check if all required fields are present
func RequiredFiledsCheck(body map[string]string, required []string) string {
	for _, field := range required {
		if len(body[field]) == 0 {
			return field
		}
	}
	return ""
}

// ReadRequestBody - read raw body from request
func ReadRequestBody(r *http.Request) (map[string]string, bool) {
	body := map[string]string{}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return body, false
	}
	defer r.Body.Close()
	json.Unmarshal(b, &body)
	return body, true
}
