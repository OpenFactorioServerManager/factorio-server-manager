package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// NilVersion represents an empty version number
var NilVersion = Version{0, 0, 0, 0}

// Version represents a semantic version and build number
type Version [4]uint

func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d.%d", v[0], v[1], v[2], v[3])
}

// MarshalText implements encoding.TextMarshaller for Version
func (v Version) MarshalText() (text []byte, err error) {
	return []byte(v.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaller for Version
func (v *Version) UnmarshalText(text []byte) error {
	parts := strings.SplitN(string(text), ".", 4)
	for i, part := range parts {
		p, err := strconv.ParseUint(part, 10, 32)
		if err != nil {
			return err
		}
		v[i] = uint(p)
	}
	return nil
}

// Equals returns true if both version are equal
func (v Version) Equals(b Version) bool {
	return v[0] == b[0] && v[1] == b[1] && v[2] == b[2] && v[3] == b[3]
}

// Less returns true if the receiver version is less than the argument version
func (v Version) Less(b Version) bool {
	switch {
	case v[0] < b[0]:
		return true
	case v[0] == b[0] && v[1] < b[1]:
		return true
	case v[0] == b[0] && v[1] == b[1] && v[2] < b[2]:
		return true
	case v[0] == b[0] && v[1] == b[1] && v[2] == b[2] && v[3] < b[3]:
		return true
	default:
		return false
	}
}

// Greater returns true if the receiver version is greater than the argument version
func (v Version) Greater(b Version) bool { return !v.Equals(b) && !v.Less(b) }

func (v Version) ge(b Version) bool { return v.Equals(b) || v.Greater(b) }
func (v Version) le(b Version) bool { return v.Equals(b) || v.Less(b) }

// Compare returns true if the comparison between the two version operands is valid.
// Supported ops are: ==, !=, >, <, >=, <=
func (v Version) Compare(b Version, op string) bool {
	switch op {
	case "==":
		return v.Equals(b)
	case "!=":
		return !v.Equals(b)
	case ">":
		return v.Greater(b)
	case "<":
		return v.Less(b)
	case ">=":
		return v.ge(b)
	case "<=":
		return v.le(b)
	default:
		panic("unsupported operator")
	}
}

// version24 is the 24-bit (8, 8, 8) version structure
type version24 Version

func (v version24) MarshalBinary() (data []byte, err error) {
	data = []byte{byte(v[0]), byte(v[1]), byte(v[2])}
	return data, nil
}

func (v *version24) UnmarshalBinary(data []byte) error {
	if len(data) < 3 {
		return errors.New("version24.UnmarshalBinary: too few bytes")
	}
	v[0] = uint(data[0])
	v[1] = uint(data[1])
	v[2] = uint(data[2])
	return nil
}

// version24 is a 48-bit (16, 16, 16) optimized version structure
type version48 Version

func (v *version48) ReadFrom(r io.Reader, game Version) error {
	for i := 0; i < 3; i++ {
		n, err := readOptimUint(r, game, 16)
		if err != nil {
			return fmt.Errorf("read part %d of version: %v", i, err)
		}
		v[i] = uint(n)
	}
	return nil
}

// version64 is the 64-bit (16, 16, 16, 16) version structure with build component.
type version64 Version

func (v version64) MarshalBinary() (data []byte, err error) {
	data = make([]byte, 8)
	binary.LittleEndian.PutUint16(data[0:2], uint16(v[0]))
	binary.LittleEndian.PutUint16(data[2:4], uint16(v[1]))
	binary.LittleEndian.PutUint16(data[4:6], uint16(v[2]))
	binary.LittleEndian.PutUint16(data[6:8], uint16(v[3]))
	return data, nil
}

func (v *version64) UnmarshalBinary(data []byte) error {
	if len(data) < 8 {
		return errors.New("version64.UnmarshalBinary: too few bytes")
	}
	v[0] = uint(binary.LittleEndian.Uint16(data[0:2]))
	v[1] = uint(binary.LittleEndian.Uint16(data[2:4]))
	v[2] = uint(binary.LittleEndian.Uint16(data[4:6]))
	v[3] = uint(binary.LittleEndian.Uint16(data[6:8]))
	return nil
}
