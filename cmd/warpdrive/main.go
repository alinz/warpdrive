package warpdrive

import "github.com/asdine/storm"
import "path/filepath"
import "os"

var (
	_bundlePath    string
	_documentPath  string
	_warpdrivePath string
	_platform      string
	_deviceCert    string
	_deviceKey     string
	_caCert        string
	_db            *storm.DB
)

type currentRelease struct {
	ID uint64 `storm:"id,increment"`
}

func dbPath() string {
	return filepath.Join(_warpdrivePath, "/warpdrive.db")
}

// Init initialize warpdrive
func Init(bundlePath, documentPath, platform, deviceCert, deviceKey, caCert string) error {
	_bundlePath = bundlePath
	_documentPath = documentPath
	_platform = platform
	_deviceCert = deviceCert
	_deviceKey = deviceKey
	_caCert = caCert

	_warpdrivePath = filepath.Join(_documentPath, "/warpdrive")
	os.MkdirAll(_warpdrivePath, os.ModePerm)

	db, err := storm.Open(dbPath())
	if err != nil {
		return err
	}

	_db = db

	return nil
}
