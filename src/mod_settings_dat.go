package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
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

//func (d *FModData) Decode(file *os.File) error {
func (d *FModData) Decode() error {
	var version version64
	var versionB [8]byte

	file, err := os.Open(filepath.Join(config.FactorioModsDir, "mod-settings.dat"))
	if err != nil {
		log.Printf("could not open mod-settings.dat")
		return err
	}

	err = binary.Read(file, binary.LittleEndian, versionB[:])
	if err != nil {
		log.Printf("could not read version: %s", err)
	}
	err = version.UnmarshalBinary(versionB[:])
	if err != nil {
		log.Printf("Error loading Version: %s", err)
		return err
	}
	d.Version = version

	d.Data, err = d.unmarshalTree(file)
	if err != nil {
		log.Printf("error loading Data: %s", err)
		return err
	}

	return nil
}

func (d *FModData) unmarshalTree(file *os.File) (interface{}, error) {
	//type of embedded data
	var _type byte
	err := binary.Read(file, binary.LittleEndian, &_type)
	if err != nil {
		log.Printf("could not read first binary: %v", err)
		return nil, err
	}

	//anyType flag ... not useful
	_, err = d.unmarshalBool(file)
	if err != nil {
		log.Printf("error loading anyType bool: %s", err)
		return nil, err
	}

	switch _type {
	case BOOL:
		return d.unmarshalBool(file)
	case DOUBLE:
		return d.unmarshalDouble(file)
	case STRING:
		return d.unmarshalString(file)
	case LIST:
		return d.unmarshalList(file)
	case DICT:
		return d.unmarshalDict(file)
	default:
		return nil, fmt.Errorf("Unknown type: %s ", err)
	}
}

func (d *FModData) unmarshalBool(file *os.File) (bool, error) {
	var _data byte
	err := binary.Read(file, binary.LittleEndian, &_data)
	if err != nil {
		log.Printf("could not read boolean byte: %s", err)
		return false, err
	}

	return _data != 0, nil
}

func (d *FModData) unmarshalDouble(file *os.File) (float64, error) {
	var _data float64
	err := binary.Read(file, binary.LittleEndian, &_data)
	if err != nil {
		log.Printf("could not read double-value: %s", err)
		return 0, err
	}

	return _data, nil
}

func (d *FModData) unmarshalList(file *os.File) ([]interface{}, error) {
	var length uint32
	length, err := readOptimUint(file, Version(d.Version), 32)
	if err != nil {
		log.Printf("could not read list length")
		return nil, err
	}

	list := make([]interface{}, length)
	for i := uint32(0); i < length; i++ {
		list[i], err = d.unmarshalTree(file)
		if err != nil {
			log.Printf("could not read tree of list: %s", err)
			return nil, err
		}
	}

	return list, nil
}

func (d *FModData) unmarshalDict(file *os.File) (map[string]interface{}, error) {
	var length uint32
	err := binary.Read(file, binary.LittleEndian, &length)
	if err != nil {
		log.Printf("could not read dict length: %s", err)
		return nil, err
	}

	dict := make(map[string]interface{})

	for i := uint32(0); i < length; i++ {
		key, err := d.unmarshalString(file)

		if err != nil {
			log.Printf("error loading key: %s", err)
			return nil, err
		}

		dict[key], err = d.unmarshalTree(file)
		if err != nil {
			log.Printf("error loading unmarshalTree: %s", err)
			return nil, err
		}
	}

	return dict, nil
}

func (d *FModData) unmarshalString(file *os.File) (string, error) {
	// read "empty" flag
	empty, err := d.unmarshalBool(file)
	if err != nil {
		log.Printf("error loading empty flag of string: %s", err)
		return "", err
	}

	if empty {
		return "", nil
	}

	key, err := readString(file, Version(d.Version), false)
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

	tree, err := d.marshalTree(d.Data)
	if err != nil {
		log.Printf("error loading first tree: %s", err)
		return nil, err
	}
	output = append(output, tree...)

	return output, nil
}

func (d *FModData) marshalTree(data interface{}) (output []byte , err error) {
	// get type
	_type := reflect.TypeOf(data).Kind()

	// write any-type flag
	anyTypeFlag := d.marshalBool(_type == reflect.String)

	var typeByte byte
	var marshalledBytes []byte

	switch _type {
	case reflect.Bool:
		typeByte = BOOL
		marshalledBytes = []byte{d.marshalBool(data.(bool))}
		break
	case reflect.Float64:
		{
			floatValue, err := d.marshalFloat64(data.(float64))
			if err != nil {
				log.Printf("could not read float64-value: %s", err)
				return nil, err
			}
			typeByte = DOUBLE
			marshalledBytes = floatValue
		}
		break
	case reflect.String:
		typeByte = STRING
		marshalledBytes = d.marshalString(data.(string))
		break
	case reflect.Array:
		// List
		list, err := d.marshalList(data.([]interface{}))
		if err != nil {
			log.Printf("could not read List: %s", err)
			return nil, err
		}
		typeByte = LIST
		marshalledBytes = list
		break
	case reflect.Map:
		// Dict
		_map, err := d.marshalDict(data.(map[string]interface{}))
		if err != nil {
			log.Printf("could not read Dict: %s", err)
			return nil, err
		}
		typeByte = DICT
		marshalledBytes = _map
		break
	}

	output = append(output, typeByte)
	output = append(output, anyTypeFlag)
	output = append(output, marshalledBytes...)

	return output, nil
}

func (d *FModData) marshalBool(data bool) byte {
	if data {
		return 0x1
	} else {
		return 0x0
	}
}

func (d *FModData) marshalFloat64(data float64) ([]byte, error) {
	var buf bytes.Buffer

	err := binary.Write(&buf, binary.LittleEndian, data)
	if err != nil {
		log.Printf("could not write data into buffer: %s", err)
		return nil, err
	}

	return buf.Bytes(), nil
}

func (d *FModData) marshalString(data string) []byte {
	var output []byte

	length := uint32(len(data))
	output = []byte{d.marshalBool(length == 0)}

	output = append(output, d.marshalOptimUint(length)...)
	stringBytes := []byte(data)
	return append(output, stringBytes...)
}

func (d *FModData) marshalOptimUint(data uint32) []byte {
	if data < 256 {
		intBinary := []byte{byte(data)}
		return intBinary[:]
	} else {
		var intBinary [4]byte
		binary.LittleEndian.PutUint32(intBinary[:], data)
		return append([]byte{0xff}, intBinary[:]...)
	}
}

func (d *FModData) marshalList(data []interface{}) ([]byte, error) {
	var output []byte

	length := uint32(len(data))
	output = d.marshalOptimUint(length)

	for i := uint32(0); i < length; i++ {
		tree, err := d.marshalTree(data[i])
		if err != nil {
			log.Printf("error loading tree of list-element: %s", err)
			return nil, err
		}
		output = append(output, tree...)
	}

	return output, nil
}

func (d *FModData) marshalDict(data map[string]interface{}) ([]byte, error) {
	var output []byte

	length := uint32(len(data))

	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], length)
	output = append(output, buf[:]...)

	for key, value := range data {
		output = append(output, d.marshalString(key)...)
		tree, err := d.marshalTree(value)
		if err != nil {
			log.Printf("error loading tree of dict-element: %s", err)
			return nil, err
		}
		output = append(output, tree...)
	}

	return output, nil
}
