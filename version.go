package warpdrive

//the following two variables will be set during the build

var (
	//VERSION describe the version number eider dev or hash value
	VERSION string
	//LONGVERSION describe the version number eider dev or hash value
	LONGVERSION string
)

func init() {
	if VERSION == "" {
		VERSION = "dev"
		LONGVERSION = "development"
	}
}
