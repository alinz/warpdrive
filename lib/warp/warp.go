package warp

import (
	"encoding/binary"
	"errors"
	"io"
	"os"

	"path/filepath"

	"fmt"

	"bytes"

	"github.com/pressly/warpdrive/lib/crypto"
)

type fileHeader struct {
	FileSize int64     // 8 bytes
	Hash     [20]byte  // 20 bytes
	Name     [996]byte // 996 bytes -> total is 1024bytes or 1kb
}

func (fh *fileHeader) read(r io.Reader) error {
	headerReader := io.LimitReader(r, 1024)
	return binary.Read(headerReader, binary.LittleEndian, fh)
}

// Warp is a type to support new file encoding
type Warp struct {
	w io.Writer
	r io.Reader
}

// AddFile adds a new filw inside warp file
func (w *Warp) AddFile(name, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	fileInfo, _ := file.Stat()
	defer file.Close()

	//calculate hash
	hash, err := crypto.Hash(file)
	if err != nil {
		return err
	}

	// we need to reset the file to 0, 0, so copy io can performe.
	file.Seek(0, 0)

	fh := fileHeader{
		FileSize: fileInfo.Size(),
	}

	nameLength := len(name)
	if nameLength > 996 {
		return errors.New("name of the file is too big")
	}

	copy(fh.Name[:], []byte(name)[:nameLength])
	copy(fh.Hash[:], hash[:20])
	err = binary.Write(w.w, binary.LittleEndian, &fh)

	if err != nil {
		return err
	}

	io.Copy(w.w, file)

	return nil
}

func (w *Warp) Extract(path string) error {
	header := &fileHeader{}

	err := header.read(w.r)
	for {
		if err != nil {
			return err
		}

		targetPath := filepath.Join(path, string(header.Name[:]))
		fmt.Println(len(string(header.Name[:])))
		// we read the header so now we need to create a folder and file under the
		// path + fileName and write the content of the file into it
		dir := filepath.Dir(targetPath)
		// filename := filepath.Base(targetPath)

		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}

		fmt.Println(targetPath, len(targetPath))

		file, err1 := os.Create(targetPath)
		if err1 != nil {
			return err1
		}

		defer file.Close()

		fileContentReader := io.LimitReader(w.r, header.FileSize)
		bytesWritten, err := io.Copy(file, fileContentReader)
		if err != nil {
			return err
		}

		if bytesWritten != header.FileSize {
			return fmt.Errorf("file size is mismatched in header for %s", targetPath)
		}

		file.Seek(0, 0)
		hash, err := crypto.Hash(file)
		if err != nil {
			return err
		}

		if !bytes.Equal(hash, header.Hash[:]) {
			return fmt.Errorf("hash mismatched for %s", targetPath)
		}

		err = header.read(w.r)
	}

	return nil
}

// NewWriter creates a new Warp for write to
func NewWriter(w io.Writer) *Warp {
	return &Warp{w: w}
}

// NewReader creates a new Warp for read from
func NewReader(r io.Reader) *Warp {
	return &Warp{r: r}
}
