package main

import (
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
