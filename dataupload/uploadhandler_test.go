package dataupload

import (
	"github.com/gaia-docker/tugbot-result-service/pool"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func TestUploadHandlerEmptyBody(t *testing.T) {

	handler := NewUploadHandler(pool.NewHub())

	assert.HTTPError(t, handler.ServeHTTP, "POST", "/results", url.Values{})
}

/*
func TestUploadHandler(t *testing.T) {

	req, err := http.NewRequest("POST", "/results", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := NewUploadHandler(pool.NewHub())
	http.Handle("/results", handler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}
*/
