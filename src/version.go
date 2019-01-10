package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// NilVersion represents an empty version number
var NilVersion = Version{0, 0, 0}

// Version represents a semantic version
type Version [3]uint

func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v[0], v[1], v[2])
}

// MarshalText implements encoding.TextMarshaller for Version
func (v Version) MarshalText() (text []byte, err error) {
	return []byte(v.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaller for Version
func (v *Version) UnmarshalText(text []byte) error {
	parts := strings.SplitN(string(text), ".", 3)
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
	return v[0] == b[0] && v[1] == b[1] && v[2] == b[2]
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

// version24 is the 48-bit (16, 16, 16) version structure
type version48 Version

func (v version48) MarshalBinary() (data []byte, err error) {
	data = make([]byte, 6)
	binary.LittleEndian.PutUint16(data[0:2], uint16(v[0]))
	binary.LittleEndian.PutUint16(data[2:4], uint16(v[1]))
	binary.LittleEndian.PutUint16(data[4:6], uint16(v[2]))
	return data, nil
}

func (v *version48) UnmarshalBinary(data []byte) error {
	if len(data) < 6 {
		return errors.New("version48.UnmarshalBinary: too few bytes")
	}
	v[0] = uint(binary.LittleEndian.Uint16(data[0:2]))
	v[1] = uint(binary.LittleEndian.Uint16(data[2:4]))
	v[2] = uint(binary.LittleEndian.Uint16(data[4:6]))
	return nil
}

// version64 is the 64-bit (16, 16, 16, 16) version structure with build component.
type version64 struct {
	Version
	Build uint
}

func (v version64) MarshalBinary() (data []byte, err error) {
	data = make([]byte, 8)
	binary.LittleEndian.PutUint16(data[0:2], uint16(v.Version[0]))
	binary.LittleEndian.PutUint16(data[2:4], uint16(v.Version[1]))
	binary.LittleEndian.PutUint16(data[4:6], uint16(v.Version[2]))
	binary.LittleEndian.PutUint16(data[6:8], uint16(v.Build))
	return data, nil
}

func (v *version64) UnmarshalBinary(data []byte) error {
	if len(data) < 8 {
		return errors.New("version64.UnmarshalBinary: too few bytes")
	}
	v.Version[0] = uint(binary.LittleEndian.Uint16(data[0:2]))
	v.Version[1] = uint(binary.LittleEndian.Uint16(data[2:4]))
	v.Version[2] = uint(binary.LittleEndian.Uint16(data[4:6]))
	v.Build = uint(binary.LittleEndian.Uint16(data[6:8]))
	return nil
}

func (v version64) String() string {
	return fmt.Sprintf("%d.%d.%d.%d", v.Version[0], v.Version[1], v.Version[2], v.Build)
}

func (v *version64) UnmarshalText(text []byte) error {
	parts := strings.SplitN(string(text), ".", 4)
	for i, part := range parts {
		p, err := strconv.ParseUint(part, 10, 32)
		if err != nil {
			return err
		}
		v.Version[i] = uint(p)
	}
	return nil
}
