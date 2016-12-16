package config

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"sort"
)

// Version is simple structure which holds information
// about which cycle uses which version
type Version struct {
	Current string   `json:"current"`
	Cached  []string `jsong:"cached"`
}

// Add adds a version to internal cache, it makes sure there are no duplicates
func (v *Version) Add(version string) {
	found := false
	for _, value := range v.Cached {
		if value == version {
			found = true
			break
		}
	}

	if !found {
		v.Cached = append(v.Cached, version)
	}
}

// SetCurrent sets a version as current for Cycle, you have an option
// to add it to cache. The reason you might not to set it to cache is
// anything added to cache means download from warpdrive server
func (v *Version) SetCurrent(version string, cache bool) {
	v.Current = version
	if cache {
		v.Add(version)
	}
}

// SortCached sort caches, it helps simplify the caller
func (v *Version) SortCached() {
	sort.Strings(v.Cached)
}

// VersionMap is a simple structure wrapping Version
// it provides basic structure to find and manipulate
// each version in an easy way
type VersionMap struct {
	values map[string]*Version
}

// Save saves the values of VersionMap to target path.
// Path must ends with filename and must be `versions.warp`
func (vm *VersionMap) Save(documentPath string) error {
	path := VersionPath(documentPath)

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer file.Close()

	return json.NewEncoder(file).Encode(vm.values)
}

// Load loads the values of VersionMap from target path.
// Path must ends with filename and must be `versions.warp`
func (vm *VersionMap) Load(documentPath string) error {
	path := VersionPath(documentPath)

	file, err := os.Open(path)
	if err != nil {
		return err
	}

	defer file.Close()

	return json.NewDecoder(file).Decode(vm.values)
}

// Version gets a verion struct based on cycle name.
// if the cycle does not exist, it will create a new one,
// assign it to map and return the value
// most of the logic here is for making it easier for caller.
func (vm *VersionMap) Version(cycle string) *Version {
	value, ok := vm.values[cycle]
	if !ok {
		value = &Version{}
		vm.values[cycle] = value
	}

	return value
}

// CurrentVersion returns the current version sets in given cycle
func (vm *VersionMap) CurrentVersion(cycle string) string {
	value := vm.Version(cycle)
	return value.Current
}

// SetCurrentVersion assing a new version to cycle. You can add it to cache as well
func (vm *VersionMap) SetCurrentVersion(cycle, version string, cache bool) {
	value := vm.Version(cycle)
	value.SetCurrent(version, cache)
}

// AddVersion adds the version to given cycle's cached and there is an option to
// set the version as current.
func (vm *VersionMap) AddVersion(cycle, version string, isCurrent bool) {
	value := vm.Version(cycle)
	if isCurrent {
		value.SetCurrent(version, true)
	} else {
		value.Add(version)
	}
}

// VersionPath returns the proper path for loading versions.warp
func VersionPath(path string) string {
	return filepath.Join(path, "versions.warp")
}

// NewVersionMapFromReader creates VersionMap from io.Reader
// added this function to simplify the creation of versionMap
func NewVersionMapFromReader(r io.Reader) (*VersionMap, error) {
	versionMap := NewVersionMap()
	err := json.NewDecoder(r).Decode(versionMap)
	return versionMap, err
}

// NewVersionMap creates a new VersionMap
func NewVersionMap() *VersionMap {
	versionMap := &VersionMap{
		values: make(map[string]*Version),
	}
	return versionMap
}
