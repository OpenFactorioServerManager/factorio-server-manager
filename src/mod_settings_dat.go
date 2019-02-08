package main

import (
	"encoding/binary"
	"io"
	"log"
)

const (
	NONE    = 0
	BOOL    = 1
	DOUBLE  = 2
	STRING  = 3
	LIST    = 4
	DICT    = 5
)

type FModData struct {
	Version version64
	Data    interface{}
}

func (d *FModData) Decode(file io.Reader) error {
	var version version64
	var versionB [8]byte

	err := binary.Read(file, binary.LittleEndian, versionB[:])
	if err != nil {
		log.Printf("could not read version: %s", err)
	}
	err = version.UnmarshalBinary(versionB[:])
	if err != nil {
		log.Printf("Error loading Version: %s", err)
		return err
	}
	d.Version = version

	d.Data, err = readTree(file)
	if err != nil {
		log.Printf("error loading Data: %s", err)
		return err
	}

	return nil
}

func readStringSettings(file io.Reader) (string, error) {
	// read "empty" flag
	empty, err := readBool(file)
	if err != nil {
		log.Printf("error loading empty flag of string: %s", err)
		return "", err
	}

	if empty {
		return "", nil
	}

	key, err := readString(file, FactorioServ.Version, false)
	if err != nil {
		log.Printf("could not read key-string: %s", err)
		return "", err
	}

	return key, nil
}

func (d *FModData) Encode() ([]byte, error) {
	var output []byte

	_bytes, err := d.Version.MarshalBinary()
	if err != nil {
		log.Printf("couldn't create binary from version: %s", err)
		return nil, err
	}

	output = append(output, _bytes...)

	tree, err := writeTree(d.Data)
	if err != nil {
		log.Printf("error loading first tree: %s", err)
		return nil, err
	}
	output = append(output, tree...)

	return output, nil
}
