package dataupload

import (
	"archive/tar"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
)

const tarFileName = "uploader_test.tar"

var testFiles = []struct {
	Name, Body string
}{
	{"readme.txt", "This archive contains some text files."},
	{"tugbot.txt", "tugbot:\ntugbot\nresult\nservice"},
	{"todo.txt", "Get animal handling licence.\nWrite more examples."},
	{"folder/", ""},
}

func TestUpload(t *testing.T) {

	createTestTar()

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

	cleanup(resultDir)
}

func TestUploadFileNotExist(t *testing.T) {

	uploader := &TarUploader{}
	reader, err := os.Open("no-file")
	resultDir, err := uploader.Upload(reader)

	assert.Error(t, err)
	cleanup(resultDir)
}

func existInTestFiles(uploadedFiles []os.FileInfo, expectedFile string) bool {

	for _, f := range uploadedFiles {
		if f.Name() == strings.TrimSuffix(expectedFile, "/") {
			return true
		}
	}

	return false
}

func cleanup(resultDir *string) {

	os.Remove(tarFileName)
	os.RemoveAll(*resultDir)
}

func createTestTar() {

	tarfile, err := os.Create(tarFileName)
	if err != nil {
		return
	}
	defer tarfile.Close()
	tw := tar.NewWriter(tarfile)
	defer tw.Close()
	// add each file as needed into the current tar archive
	for _, file := range testFiles {
		var typeFlag = byte(tar.TypeReg)
		if strings.HasSuffix(file.Name, "/") {
			typeFlag = tar.TypeDir
		}
		hdr := &tar.Header{
			Name:     file.Name,
			Mode:     0600,
			Size:     int64(len(file.Body)),
			Typeflag: typeFlag,
		}
		if err := tw.WriteHeader(hdr); err != nil {
			log.Fatalln(err)
		}
		if _, err := tw.Write([]byte(file.Body)); err != nil {
			log.Fatalln(err)
		}
	}
}
