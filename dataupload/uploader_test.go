package dataupload

import (
	"fmt"
	"github.com/gaia-docker/tugbot-result-service/testutils"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

const tarFileName = "uploader_test.tar"

var testFiles = []testutils.FileDescriptor{
	{"readme.txt", "This archive contains some text files."},
	{"tugbot.txt", "tugbot:\ntugbot\nresult\nservice"},
	{"todo.txt", "Get animal handling licence.\nWrite more examples."},
	{"folder/", ""},
}

func TestUpload(t *testing.T) {

	testutils.CreateTarWithName(tarFileName, testFiles)

	uploader := &TarUploader{}
	reader, err := os.Open(tarFileName)
	resultDir, err := uploader.Upload(reader)

	assert.NoError(t, err)
	_, err = os.Stat(*resultDir)
	assert.NoError(t, err)
	files, err := ioutil.ReadDir(*resultDir)
	assert.NoError(t, err)
	// check that every test file was uploaded
	for _, f := range testFiles {
		assert.True(t, existInTestFiles(files, f.Name), fmt.Sprintf("%s should exist in tar file", f.Name))
	}

	cleanup(resultDir, tarFileName)
}

func TestUploadFileNotExist(t *testing.T) {

	uploader := &TarUploader{}
	reader, err := os.Open("no-file")
	resultDir, err := uploader.Upload(reader)

	assert.Error(t, err)
	cleanup(resultDir, "no-file")
}

func existInTestFiles(uploadedFiles []os.FileInfo, expectedFile string) bool {

	for _, f := range uploadedFiles {
		if f.Name() == strings.TrimSuffix(expectedFile, "/") {
			return true
		}
	}

	return false
}

func cleanup(resultDir *string, tarFileName string) {

	os.Remove(tarFileName)
	if resultDir != nil {
		os.RemoveAll(*resultDir)
	}
}
