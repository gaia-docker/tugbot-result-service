package dataupload

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/gaia-docker/tugbot-result-service/pool"
	"io/ioutil"
	"net/http"
)

// UploadHandler responds to /results http request, which is the result-service rest API for uploading results
type UploadHandler struct {
	hub *pool.Hub
}

// NewUploadHandler creates UploadHandler instance
func NewUploadHandler(hub *pool.Hub) *UploadHandler {

	return &UploadHandler{
		hub: hub,
	}
}

func (uh *UploadHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	retStatus := http.StatusOK
	body, err := getBody(request)
	if err != nil {
		retStatus = http.StatusBadRequest
		log.Error("Error fetching request body. ", err)
	} else {
		uh.hub.Broadcast(body)
	}
	writer.WriteHeader(retStatus)
}

func getBody(request *http.Request) (*string, error) {

	requestBody := request.Body
	if requestBody == nil {
		return nil, errors.New("Empty request body")
	}
	body, err := ioutil.ReadAll(requestBody)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	ret := string(body)

	return &ret, nil
}
