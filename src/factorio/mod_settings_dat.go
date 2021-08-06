package factorio

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
)

const (
	NONE   = 0
	BOOL   = 1
	DOUBLE = 2
	STRING = 3
	LIST   = 4
	DICT   = 5
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

	if Version(version).Greater(Version{0, 17, 0, 0}) {
		//FIXME correct naming
		var b [1]byte
		_, err = file.Read(b[:])
		if err != nil {
			return fmt.Errorf("read first random 0.17 byte: %v", err)
		}
	}

	d.Data, err = readTree(file, Version(d.Version))
	if err != nil {
		log.Printf("error loading Data: %s", err)
		return err
	}

	return nil
}

func (d *FModData) Encode() ([]byte, error) {
	var output []byte

	_bytes, err := d.Version.MarshalBinary()
	if err != nil {
		log.Printf("couldn't create binary from version: %s", err)
		return nil, err
	}

	output = append(output, _bytes...)

	if Version(d.Version).Greater(Version{0, 17, 0, 0}) {
		output = append(output, byte(0))
	}

	tree, err := writeTree(d.Data)
	if err != nil {
		log.Printf("error loading first tree: %s", err)
		return nil, err
	}
	output = append(output, tree...)

	return output, nil
}
