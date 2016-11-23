package data

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// Version holds version of the upgradeable application in uint64
type Version uint64

const MaxMajor = 0xffff000000000000
const MaxMinor = 0x0000ffff00000000
const MaxPatch = 0x00000000ffffffff

func ParseVersion(version string) (Version, error) {
	decode, err := VersionDecode(version)
	if err != nil {
		return 0, err
	}

	var value uint64

	decode[0] = (decode[0] << 48) & MaxMajor
	decode[1] = (decode[1] << 32) & MaxMinor
	decode[2] = decode[2] & MaxPatch

	value = decode[0] | decode[1] | decode[2]

	return Version(value), nil
}

func VersionToInt(version Version) uint64 {
	return uint64(version)
}

// VersionEncode convert Version to string represantation
func VersionEncode(major, minor, patch uint64) string {
	return strconv.FormatUint(major, 10) + "." +
		strconv.FormatUint(minor, 10) + "." +
		strconv.FormatUint(patch, 10)
}

// VersionDecode gets an major.minor.patch and returns array of parsed version
func VersionDecode(value string) ([]uint64, error) {
	versionSegments := strings.Split(value, ".")

	if len(versionSegments) != 3 {
		return nil, fmt.Errorf("Version should have 3 parts, got %d", len(versionSegments))
	}

	major, err := strconv.ParseUint(versionSegments[0], 10, 16)
	if err != nil {
		return nil, fmt.Errorf("Major part of Version is not parrsable")
	}

	minor, err := strconv.ParseUint(versionSegments[1], 10, 16)
	if err != nil {
		return nil, fmt.Errorf("Minor part of Version is not parrsable")
	}

	patch, err := strconv.ParseUint(versionSegments[2], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("Patch part of Version is not parrsable")
	}

	return []uint64{major, minor, patch}, nil
}

// MarshalJSON for type Version
func (v Version) MarshalJSON() ([]byte, error) {
	value := uint64(v)

	major := (value & MaxMajor) >> 48
	minor := (value & MaxMinor) >> 32
	patch := (value & MaxPatch)

	version := VersionEncode(major, minor, patch)

	return json.Marshal(version)
}

// UnmarshalJSON for type Version
func (v *Version) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("Version should be a string, got %s", data)
	}

	version, err := VersionDecode(s)
	if err != nil {
		return err
	}

	var value uint64

	version[0] = (version[0] << 48) & MaxMajor
	version[1] = (version[1] << 32) & MaxMinor
	version[2] = version[2] & MaxPatch

	value = version[0] | version[1] | version[2]

	*v = Version(value)
	return nil
}

func (v Version) ValueAsInt() uint64 {
	return uint64(v)
}

func VersionAdd(version Version, major, minor, patch uint64) Version {
	value := uint64(version)

	major = major + ((value & MaxMajor) >> 48)
	minor = minor + ((value & MaxMinor) >> 32)
	patch = patch + (value & MaxPatch)

	major = (major << 48) & MaxMajor
	minor = (minor << 32) & MaxMinor
	patch = patch & MaxPatch

	value = major | minor | patch

	return Version(value)
}

func MaskVersion(version Version, major, minor, patch uint64) Version {
	value := uint64(version)

	v1 := (value & MaxMajor) >> 48
	v2 := (value & MaxMinor) >> 32
	v3 := value & MaxPatch

	if major == 0 {
		v1 = 0
	}
	if minor == 0 {
		v2 = 0
	}
	if patch == 0 {
		v3 = 0
	}

	v1 = (v1 << 48) & MaxMajor
	v2 = (v2 << 32) & MaxMinor
	v3 = v3 & MaxPatch

	value = v1 | v2 | v3

	return Version(value)
}
