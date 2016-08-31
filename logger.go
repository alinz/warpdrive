package warpdrive

import (
	"log"
	"os"
)

//Logger global logger
var Logger *log.Logger

func init() {
	Logger = log.New(os.Stdout, "logger: ", log.Lshortfile)
}
