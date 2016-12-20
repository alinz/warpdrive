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
	Current   string          `json:"current"`
	Available map[string]bool `json:"available"`
}

// Add adds a version to internal cache, it makes sure there are no duplicates
func (v *Version) Add(version string, isBundle bool) {
	if _, ok := v.Available[version]; !ok {
		v.Available[version] = isBundle
	}
}

// SetCurrent sets a version as current for Cycle, you have an option
// to add it to cache. The reason you might not to set it to cache is
// anything added to cache means download from warpdrive server
func (v *Version) SetCurrent(version string, isBundle bool) {
	v.Current = version
	v.Add(version, isBundle)
}

// SortAvailable sort caches, it helps simplify the caller
func (v *Version) SortAvailable() []string {
	var versions []string
	for version := range v.Available {
		versions = append(versions, version)
	}
	sort.Strings(versions)
	return versions
}

// VersionMap is a simple structure wrapping Version
// it provides basic structure to find and manipulate
// each version in an easy way
type VersionMap struct {
	ActiveCycle string              `json:"active_cycle"`
	Cycles      map[string]*Version `json:"cycles"`
}

// Save saves the VersionMap to target path.
// Path must ends with filename and must be `versions.warp`
func (vm *VersionMap) Save(documentPath string) error {
	path := VersionPath(documentPath)

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer file.Close()

	return json.NewEncoder(file).Encode(vm)
}

// Load loads the Cycles of VersionMap from target path.
// Path must ends with filename and must be `versions.warp`
func (vm *VersionMap) Load(documentPath string) error {
	path := VersionPath(documentPath)

	file, err := os.Open(path)
	if err != nil {
		return err
	}

	defer file.Close()

	return json.NewDecoder(file).Decode(vm)
}

// Version gets a verion struct based on cycle name.
// if the cycle does not exist, it will create a new one,
// assign it to map and return the value
// most of the logic here is for making it easier for caller.
func (vm *VersionMap) Version(cycle string) *Version {
	if vm.Cycles == nil {
		vm.Cycles = make(map[string]*Version)
	}

	value, ok := vm.Cycles[cycle]
	if !ok {
		value = &Version{
			Available: make(map[string]bool),
		}
		vm.Cycles[cycle] = value
	}

	return value
}

// SetActiveCycle set the ActiveCycle to load the value during bootup
func (vm *VersionMap) SetActiveCycle(cycleName string) {
	vm.ActiveCycle = cycleName
}

// CurrentVersion returns the current version sets in given cycle
func (vm *VersionMap) CurrentVersion(cycle string) string {
	value := vm.Version(cycle)
	return value.Current
}

// SetCurrentVersion assing a new version to cycle. You can add it to cache as well
func (vm *VersionMap) SetCurrentVersion(cycle, version string, isBundle bool) {
	value := vm.Version(cycle)
	value.SetCurrent(version, isBundle)
}

// VersionPath returns the proper path for loading versions.warp
func VersionPath(documentPath string) string {
	return filepath.Join(documentPath, "warpdrive/versions.warp")
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
		Cycles: make(map[string]*Version),
	}
	return versionMap
}
