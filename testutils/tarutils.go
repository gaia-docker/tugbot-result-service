package testutils

import (
	"archive/tar"
	"log"
	"os"
	"strings"
)

type FileDescriptor struct {
	Name string
	Body string
}

func CreateTar(fileWriter *tar.Writer, testFiles []FileDescriptor) {

	defer fileWriter.Close()
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
		if err := fileWriter.WriteHeader(hdr); err != nil {
			log.Fatalln(err)
		}
		if _, err := fileWriter.Write([]byte(file.Body)); err != nil {
			log.Fatalln(err)
		}
	}
}

func CreateTarWithName(tarFileName string, testFiles []FileDescriptor) {

	tarfile, err := os.Create(tarFileName)
	if err != nil {
		return
	}
	defer tarfile.Close()
	CreateTar(tar.NewWriter(tarfile), testFiles)
}
