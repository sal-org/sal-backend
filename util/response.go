package util

import (
	"encoding/json"
	"net/http"
	CONSTANT "salbackend/constant"
)

// SetReponse - set request response with status, message
func SetReponse(w http.ResponseWriter, status string, msg string, msgType string, response map[string]interface{}) {
	w.Header().Set("Status", "200")
	response["meta"] = setMeta(status, msg, msgType)
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(response)
}

func setMeta(status string, msg string, msgType string) map[string]string {
	if len(msg) == 0 {
		if status == CONSTANT.StatusCodeBadRequest {
			msg = "Bad Request"
		} else if status == CONSTANT.StatusCodeServerError {
			msg = "Server Error"
		}
	}
	return map[string]string{
		"status":       status,
		"message":      msg,
		"message_type": msgType,
	}
}

func getHTTPStatusCode(code string) int {
	switch code {
	case CONSTANT.StatusCodeOk:
		return http.StatusOK
	case CONSTANT.StatusCodeCreated:
		return http.StatusCreated
	case CONSTANT.StatusCodeBadRequest:
		return http.StatusBadRequest
	case CONSTANT.StatusCodeServerError:
		return http.StatusInternalServerError
	}
	return http.StatusOK
}
