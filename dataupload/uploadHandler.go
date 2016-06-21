package dataupload

import (
	"compress/gzip"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gaia-docker/tugbot-result-service/pool"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
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
	params := request.URL.Query()
	mainFileName := params.Get("mainfile")
	if err != nil {
		retStatus = http.StatusBadRequest
		log.Error("Error fetching request body. ", err)
	} else {
		resultDir, err := uh.uploader.Upload(body)
		if err == nil {
			uh.broadcast(*resultDir, mainFileName)
		}
	}
	writer.WriteHeader(retStatus)
}

func (uh UploadHandler) broadcast(resultDir, mainFileName string) {

	var message string = ""
	if mainFileName != "" {
		content, err := ioutil.ReadFile(filepath.Join(resultDir, mainFileName))
		if err != nil {
			log.Errorf("error reading file content <%s> %s", mainFileName, err)
		} else {
			message = string(content)
		}
	}
	if message == "" {
		message = fmt.Sprintf("Received test results, no main file was specified <result dir=%s>", resultDir)
	}
	uh.hub.Broadcast(&message)
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
