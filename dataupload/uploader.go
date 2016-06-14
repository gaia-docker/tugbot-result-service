package dataupload

import (
	"archive/zip"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Uploader interface
type Uploader interface {
	Upload(filename string)
}

// ZipUploader implements the Uploader interface
type ZipUploader struct {
}

func (zu *ZipUploader) Upload(filename string) (*string, error) {

	reader, err := zip.OpenReader(filename)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer reader.Close()
	resultDirName := fmt.Sprintf("resultService_%s", strings.TrimSuffix(filename, ".zip"))
	os.Mkdir(resultDirName, os.ModeDir)

	for _, f := range reader.Reader.File {

		zipped, err := f.Open()
		if err != nil {
			log.Error(err)
			continue
		}
		zu.unzip(zipped, f, resultDirName)
	}

	return &resultDirName, err
}

func (zu *ZipUploader) unzip(zipped io.ReadCloser, file *zip.File, resultDirName string) {

	defer zipped.Close()

	// get the individual file name and extract the current directory
	path := filepath.Join("./", resultDirName, file.Name)

	if file.FileInfo().IsDir() {
		os.MkdirAll(path, file.Mode())
		log.Infof("Creating directory %s", path)
	} else {
		writer, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, file.Mode())
		if err != nil {
			log.Error(err)
			return
		}
		defer writer.Close()
		if _, err = io.Copy(writer, zipped); err != nil {
			log.Error(err)
			return
		}
		log.Infof("Decompressing: %s", path)
	}
}
