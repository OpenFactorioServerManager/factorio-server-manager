package factorio

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"reflect"
)

///////////////////
// Reading ////////
///////////////////
func readOptimUint(r io.Reader, v Version, bitSize int) (uint32, error) {
	var b [4]byte
	if !v.Less(Version{0, 14, 14, 0}) {
		_, err := r.Read(b[:1])
		if err != nil {
			return 0, err
		}
		if b[0] != 0xFF {
			return uint32(b[0]), nil
		}
	}

	if bitSize < 0 || bitSize > 64 || (bitSize%8 != 0) {
		panic("invalid bit size")
	}

	_, err := r.Read(b[:bitSize/8])
	if err != nil {
		return 0, err
	}

	switch bitSize {
	case 16:
		return uint32(binary.LittleEndian.Uint16(b[:2])), nil
	case 32:
		return binary.LittleEndian.Uint32(b[:4]), nil
	default:
		panic("invalid bit size")
	}
}

func readString(r io.Reader, version Version, forceOptimized bool) (s string, err error) {
	var n uint32

	if !version.Less(Version{0, 16, 0, 0}) || forceOptimized {
		n, err = readOptimUint(r, version, 32)
		if err != nil {
			return "", err
		}
	} else {
		var b [4]byte
		_, err := r.Read(b[:])
		if err != nil {
			return "", fmt.Errorf("failed to read string length: %v", err)
		}
		n = uint32(binary.LittleEndian.Uint32(b[:]))
	}

	d := make([]byte, n)
	_, err = r.Read(d)
	if err != nil {
		return "", fmt.Errorf("failed to read string: %v", err)
	}

	return string(d), nil
}

func readBool(file io.Reader) (bool, error) {
	var _data byte
	err := binary.Read(file, binary.LittleEndian, &_data)
	if err != nil {
		log.Printf("could not read boolean byte: %s", err)
		return false, err
	}

	return _data != 0, nil
}

func readDouble(file io.Reader) (float64, error) {
	var _data float64
	err := binary.Read(file, binary.LittleEndian, &_data)
	if err != nil {
		log.Printf("could not read double-value: %s", err)
		return 0, err
	}

	return _data, nil
}

func readList(file io.Reader, version Version) ([]interface{}, error) {
	var length uint32
	length, err := readOptimUint(file, version, 32)
	if err != nil {
		log.Printf("could not read list length")
		return nil, err
	}

	list := make([]interface{}, length)
	for i := uint32(0); i < length; i++ {
		list[i], err = readTree(file, version)
		if err != nil {
			log.Printf("could not read tree of list: %s", err)
			return nil, err
		}
	}

	return list, nil
}

func readDict(file io.Reader, version Version) (map[string]interface{}, error) {
	var length uint32
	err := binary.Read(file, binary.LittleEndian, &length)
	if err != nil {
		log.Printf("could not read dict length: %s", err)
		return nil, err
	}

	dict := make(map[string]interface{})

	for i := uint32(0); i < length; i++ {
		key, err := readStringSettings(file, version)

		if err != nil {
			log.Printf("error loading key: %s", err)
			return dict, err
		}

		dict[key], err = readTree(file, version)
		if err != nil {
			log.Printf("error loading readTree: %s", err)
			return dict, err
		}
	}

	return dict, nil
}

func readTree(file io.Reader, version Version) (interface{}, error) {
	//type of embedded data
	var _type byte
	err := binary.Read(file, binary.LittleEndian, &_type)
	if err != nil {
		log.Printf("could not read first binary: %v", err)
		return nil, err
	}

	//anyType flag ... not useful
	_, err = readBool(file)
	if err != nil {
		log.Printf("error loading anyType bool: %s", err)
		return nil, err
	}

	switch _type {
	case BOOL:
		return readBool(file)
	case DOUBLE:
		return readDouble(file)
	case STRING:
		return readStringSettings(file, version)
	case LIST:
		return readList(file, version)
	case DICT:
		return readDict(file, version)
	default:
		return nil, fmt.Errorf("Unknown type: %s ", err)
	}
}

///////////////////
// Writing ////////
///////////////////
func writeOptimUint(data uint32) []byte {
	if data < 256 {
		intBinary := []byte{byte(data)}
		return intBinary[:]
	} else {
		var intBinary [4]byte
		binary.LittleEndian.PutUint32(intBinary[:], data)
		return append([]byte{0xff}, intBinary[:]...)
	}
}

func writeFloat64(data float64) ([]byte, error) {
	var buf bytes.Buffer

	err := binary.Write(&buf, binary.LittleEndian, data)
	if err != nil {
		log.Printf("could not write data into buffer: %s", err)
		return nil, err
	}

	return buf.Bytes(), nil
}

func writeBool(data bool) byte {
	if data {
		return 0x1
	} else {
		return 0x0
	}
}

func writeString(data string) []byte {
	var output []byte

	length := uint32(len(data))
	// True if the string is empty ... not used by factorio, so set to false
	//output = []byte{writeBool(length == 0)}
	output = []byte{writeBool(false)}

	output = append(output, writeOptimUint(length)...)

	if length != 0 {
		stringBytes := []byte(data)
		output = append(output, stringBytes...)
	}
	return output
}

func writeList(data []interface{}) ([]byte, error) {
	var output []byte

	length := uint32(len(data))
	output = writeOptimUint(length)

	for i := uint32(0); i < length; i++ {
		tree, err := writeTree(data[i])
		if err != nil {
			log.Printf("error loading tree of list-element: %s", err)
			return nil, err
		}
		output = append(output, tree...)
	}

	return output, nil
}

func writeDict(data map[string]interface{}) ([]byte, error) {
	var output []byte

	length := uint32(len(data))

	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], length)
	output = append(output, buf[:]...)

	for key, value := range data {
		output = append(output, writeString(key)...)
		tree, err := writeTree(value)
		if err != nil {
			log.Printf("error loading tree of dict-element: %s", err)
			return nil, err
		}
		output = append(output, tree...)
	}

	return output, nil
}

func writeTree(data interface{}) (output []byte, err error) {
	// get type
	_type := reflect.TypeOf(data).Kind()

	// write any-type flag -- Not used by factorio ... so set to false
	//anyTypeFlag := writeBool(_type == reflect.String)
	anyTypeFlag := writeBool(false)

	var typeByte byte
	var marshalledBytes []byte

	switch _type {
	case reflect.Bool:
		typeByte = BOOL
		marshalledBytes = []byte{writeBool(data.(bool))}
	case reflect.Int:
		floatValue, err := writeFloat64(float64(data.(int)))
		if err != nil {
			log.Printf("could not write int to float64-value: %s", err)
			return nil, err
		}
		typeByte = DOUBLE
		marshalledBytes = floatValue
	case reflect.Int8:
		floatValue, err := writeFloat64(float64(data.(int8)))
		if err != nil {
			log.Printf("could not write int8 to float64-value: %s", err)
			return nil, err
		}
		typeByte = DOUBLE
		marshalledBytes = floatValue
	case reflect.Int16:
		floatValue, err := writeFloat64(float64(data.(int16)))
		if err != nil {
			log.Printf("could not write int16 to float64-value: %s", err)
			return nil, err
		}
		typeByte = DOUBLE
		marshalledBytes = floatValue
	case reflect.Int32:
		floatValue, err := writeFloat64(float64(data.(int32)))
		if err != nil {
			log.Printf("could not write int32 to float64-value: %s", err)
			return nil, err
		}
		typeByte = DOUBLE
		marshalledBytes = floatValue
	case reflect.Int64:
		floatValue, err := writeFloat64(float64(data.(int64)))
		if err != nil {
			log.Printf("could not write int64 to float64-value: %s", err)
			return nil, err
		}
		typeByte = DOUBLE
		marshalledBytes = floatValue
	case reflect.Uint:
		floatValue, err := writeFloat64(float64(data.(uint)))
		if err != nil {
			log.Printf("could not write uint to float64-value: %s", err)
			return nil, err
		}
		typeByte = DOUBLE
		marshalledBytes = floatValue
	case reflect.Uint8:
		floatValue, err := writeFloat64(float64(data.(uint8)))
		if err != nil {
			log.Printf("could not write uint8 to float64-value: %s", err)
			return nil, err
		}
		typeByte = DOUBLE
		marshalledBytes = floatValue
	case reflect.Uint16:
		floatValue, err := writeFloat64(float64(data.(uint16)))
		if err != nil {
			log.Printf("could not write uint16 to float64-value: %s", err)
			return nil, err
		}
		typeByte = DOUBLE
		marshalledBytes = floatValue
	case reflect.Uint32:
		floatValue, err := writeFloat64(float64(data.(uint32)))
		if err != nil {
			log.Printf("could not write uint32 to float64-value: %s", err)
			return nil, err
		}
		typeByte = DOUBLE
		marshalledBytes = floatValue
	case reflect.Uint64:
		floatValue, err := writeFloat64(float64(data.(uint64)))
		if err != nil {
			log.Printf("could not write uint64 to float64-value: %s", err)
			return nil, err
		}
		typeByte = DOUBLE
		marshalledBytes = floatValue
	case reflect.Float32:
		floatValue, err := writeFloat64(float64(data.(float32)))
		if err != nil {
			log.Printf("could not write float32 to float64-value: %s", err)
			return nil, err
		}
		typeByte = DOUBLE
		marshalledBytes = floatValue
	case reflect.Float64:
		floatValue, err := writeFloat64(data.(float64))
		if err != nil {
			log.Printf("could not write float64 to float64-value: %s", err)
			return nil, err
		}
		typeByte = DOUBLE
		marshalledBytes = floatValue
	case reflect.String:
		typeByte = STRING
		marshalledBytes = writeString(data.(string))
	case reflect.Array:
		// List
		list, err := writeList(data.([]interface{}))
		if err != nil {
			log.Printf("could not read List: %s", err)
			return nil, err
		}
		typeByte = LIST
		marshalledBytes = list
	case reflect.Map:
		// Dict
		_map, err := writeDict(data.(map[string]interface{}))
		if err != nil {
			log.Printf("could not read Dict: %s", err)
			return nil, err
		}
		typeByte = DICT
		marshalledBytes = _map
	default:
		log.Println("Unknown Datatype")
		return output, fmt.Errorf("unknown datatype")
	}

	output = append(output, typeByte)
	output = append(output, anyTypeFlag)
	output = append(output, marshalledBytes...)

	return output, nil
}
