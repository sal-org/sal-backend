package util

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
)

func GetNumberOfPages(count string, countPerPage int) int {
	ctn, _ := strconv.Atoi(count)
	return int(math.Ceil(float64(ctn) / float64(countPerPage)))
}

// GetPageNumber - get page number for pagination
func GetPageNumber(pageStr string) int {
	page, _ := strconv.Atoi(pageStr)
	if page <= 0 {
		page = 1
	}
	return page
}

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

// ReadRequestBodyInListMap - read raw body from request to list of maps
func ReadRequestBodyInListMap(r *http.Request) ([]map[string]string, bool) {
	body := []map[string]string{}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return body, false
	}
	defer r.Body.Close()
	json.Unmarshal(b, &body)
	return body, true
}
