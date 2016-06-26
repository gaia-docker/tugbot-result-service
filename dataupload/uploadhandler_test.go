package dataupload

import (
	"archive/tar"
	"compress/gzip"
	"github.com/gaia-docker/tugbot-result-service/pool"
	"github.com/gaia-docker/tugbot-result-service/testutils"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"
)

func TestUploadHandlerEmptyBody(t *testing.T) {

	handler := NewUploadHandler(pool.NewHub())

	assert.HTTPError(t, handler.ServeHTTP, "POST", "/results", url.Values{})
}

func TestUploadHandler(t *testing.T) {

	testFile := "test.tar.gz"
	createTestGzip(testFile)
	hub := pool.NewHub()
	go hub.Run()

	code := uploadResults(t, hub, testFile)
	assert.Equal(t, http.StatusOK, code)

	matches, _ := filepath.Glob("resultService*")
	cleanup(&matches[0], testFile)
	hub.CloseBroadcastChannel()
}

func uploadResults(t *testing.T, hub *pool.Hub, testFile string) int {

	file, err := os.Open(testFile)
	assert.NoError(t, err)
	defer file.Close()

	req, err := http.NewRequest("POST", "/results?mainfile=tugbot.txt", file)
	assert.NoError(t, err)

	req.Header.Add("Content-Type", "gzip")

	rr := httptest.NewRecorder()
	handler := NewUploadHandler(hub)
	handler.ServeHTTP(rr, req)

	return rr.Code
}

func createTestGzip(testFile string) {

	file, err := os.Create(testFile)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	// set up the gzip writer
	gw := gzip.NewWriter(file)
	defer gw.Close()

	testutils.CreateTar(tar.NewWriter(gw), testFiles)
}
