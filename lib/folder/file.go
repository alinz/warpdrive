package folder

import (
	"io"
	"os"
)

// PathExists checks whether path exists or not
// it sounds simple but it turns out, it requires a little bit more work
// that's why I have created this function
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// ReadFromFile a simple warpper for reading a file once the file is open.
// it closes the file resource automatically
func ReadFromFile(path string, callback func(io.Reader, error)) {
	file, err := os.Open(path)
	if err != nil {
		callback(nil, err)
		return
	}

	defer file.Close()

	callback(file, nil)
}

// WriteToFile writes data stream into file
func WriteToFile(path string, data io.Reader) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer file.Close()

	io.Copy(file, data)

	return nil
}
