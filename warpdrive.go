package warpdrive

// The following two variables will be set during the build

var (
	// VERSION describes the version number - either dev or hash value
	VERSION string

	// LONGVERSION describes the version number - either dev or hash value
	LONGVERSION string
)

func init() {
	if VERSION == "" {
		VERSION = "dev"
		LONGVERSION = "development"
	}
}
