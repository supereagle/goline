package http

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/supereagle/jenkins-pipeline/utils/json"
)

type ResponseEntity struct {
	Code       int         `json:"code"`
	Status     string      `json:"status"`
	ErrorMsg   string      `json:"error,omitempty"`
	JsonObject interface{} `json:"json_object,omitempty"`
}

func WriteResponse(resp http.ResponseWriter, code int, jsonObject interface{}, err error) {
	respEntity := ResponseEntity{
		Code:   code,
		Status: http.StatusText(code),
	}

	if jsonObject != nil {
		respEntity.JsonObject = jsonObject
	}

	if err != nil {
		respEntity.ErrorMsg = err.Error()
	}

	respStr, err := json.Marshal2JsonStr(respEntity)
	if err != nil {
		log.Errorf("Fail to marshal the response entity: %v", respEntity)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(code)
	resp.Write([]byte(respStr))
}
