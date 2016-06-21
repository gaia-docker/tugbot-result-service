package dataupload

import (
	"compress/gzip"
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/gaia-docker/tugbot-result-service/pool"
	"io"
	"net/http"
	"strings"
)

// UploadHandler responds to /results http request, which is the result-service rest API for uploading results
type UploadHandler struct {
	hub      *pool.Hub
	uploader Uploader
}

// NewUploadHandler creates UploadHandler instance
func NewUploadHandler(hub *pool.Hub) *UploadHandler {

	return &UploadHandler{
		hub:      hub,
		uploader: TarUploader{},
	}
}

func (uh UploadHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	retStatus := http.StatusOK
	body, err := getBody(request)
	if err != nil {
		retStatus = http.StatusBadRequest
		log.Error("Error fetching request body. ", err)
	} else {
		result, err := uh.uploader.Upload(body)
		if err == nil {
			uh.hub.Broadcast(result)
		}
	}
	writer.WriteHeader(retStatus)
}

func getBody(request *http.Request) (io.ReadCloser, error) {

	ret := request.Body
	if ret == nil {
		return nil, errors.New("Empty request body")
	}
	var err error
	if strings.Contains(request.Header.Get("Content-Type"), "gzip") {
		log.Debug("Recieved gzip file as body, creating gzip reader... ")
		ret, err = gzip.NewReader(ret)
	}

	return ret, err
}
