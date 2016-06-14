package dataupload

import (
	"archive/tar"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"time"
)

// Uploader interface
type Uploader interface {
	Upload(fileReader io.ReadCloser) (*string, error)
}

// TarUploader implements the Uploader interface
type TarUploader struct {
}

func (zu TarUploader) Upload(fileReader io.ReadCloser) (*string, error) {

	defer fileReader.Close()
	tarBallReader := tar.NewReader(fileReader)
	resultDirName := fmt.Sprintf("resultService_%d", time.Now().Unix())
	os.Mkdir(resultDirName, os.ModeDir)
	var err error
	var header *tar.Header

	log.Infof("Uploading results to %s", resultDirName)
	for err == nil {
		header, err = tarBallReader.Next()
		if err == nil {
			err = zu.untar(tarBallReader, header, resultDirName)
		}
	}
	if hasError(err) {
		log.Errorf("Error ocured during untar: %s", err)
		return &resultDirName, err
	}

	return &resultDirName, nil
}

func (zu *TarUploader) untar(tarBallReader *tar.Reader, header *tar.Header, resultDirName string) error {

	var ret error
	// get the individual file name and extract the current directory
	path := filepath.Join("./", resultDirName, header.Name)

	switch header.Typeflag {
	case tar.TypeDir:
		// handle directory
		log.Infof("Creating directory %s", path)
		ret = os.MkdirAll(path, os.FileMode(header.Mode))

	case tar.TypeReg, tar.TypeRegA:
		// handle normal file
		log.Infof("Untarring: %s", path)
		writer, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, header.FileInfo().Mode())
		defer writer.Close()

		if err == nil {
			_, err = io.Copy(writer, tarBallReader)
		}
		ret = err
	default:
		log.Warnf("Unable to untar type: %c in file %s", header.Typeflag, header.Name)
	}

	return ret
}

func hasError(err error) bool {

	return err != nil && err != io.EOF
}
