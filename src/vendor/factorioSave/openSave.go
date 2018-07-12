package factorioSave

import (
	"archive/zip"
	"log"
	"io"
	"errors"
)

var ErrorLevelDatNotFound = errors.New("couldn't find level.dat")

func openSave(filePath string) (io.ReadCloser, error) {
	var err error

	saveFile, err := zip.OpenReader(filePath)
	if err != nil {
		log.Printf("error opening saveFile: %s", err)
		return nil, err
	}

	for _, singleFile := range saveFile.File {
		if singleFile.FileInfo().Name() == "level.dat" {
			//open level.dat
			rc, err := singleFile.Open()
			if err != nil {
				log.Printf("Couldn't open level.dat: %s", err)
				return nil, err
			}

			return rc, nil
		}
	}

	return nil, ErrorLevelDatNotFound
}
