package dataupload

import (
	"archive/zip"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

const zipFileName = "uploader_test.zip"

var testFiles = []struct {
	Name, Body string
}{
	{"readme.txt", "This archive contains some text files."},
	{"tugbot.txt", "tugbot:\tugbot\result\nservice"},
	{"todo.txt", "Get animal handling licence.\nWrite more examples."},
}

func TestUpload(t *testing.T) {

	createTestZip()

	uploader := &ZipUploader{}
	resultDir, err := uploader.Upload(zipFileName)

	assert.NoError(t, err)
	_, err = os.Stat(*resultDir)
	assert.NoError(t, err)
	files, err := ioutil.ReadDir(*resultDir)
	assert.NoError(t, err)
	// check that every test file was uploaded
	for _, f := range testFiles {
		assert.True(t, existInTestFiles(files, f.Name))
	}

	cleanup(resultDir)
}

func existInTestFiles(uploadedFiles []os.FileInfo, expectedFile string) bool {

	for _, f := range uploadedFiles {
		if f.Name() == expectedFile {
			return true
		}
	}

	return false
}

func cleanup(resultDir *string) {

	os.Remove(zipFileName)
	os.RemoveAll(*resultDir)
}

func createTestZip() {

	zipfile, err := os.Create(zipFileName)
	if err != nil {
		return
	}
	defer zipfile.Close()
	// Create a new zip archive.
	archive := zip.NewWriter(zipfile)
	defer archive.Close()
	// Add some files to the archive.
	for _, file := range testFiles {
		f, err := archive.Create(file.Name)
		if err != nil {
			log.Fatal(err)
		}
		_, err = f.Write([]byte(file.Body))
		if err != nil {
			log.Fatal(err)
		}
	}
}