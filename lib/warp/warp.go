package warp

import (
	"encoding/binary"
	"errors"
	"io"
	"os"

	"github.com/pressly/warpdrive/lib/crypto"
)

type fileHeader struct {
	FileSize int64     //8 bytes
	Hash     [20]byte  //20 bytes
	Name     [996]byte //996 bytes -> total is 1024bytes or 1kb
}

//Warp is a type to support new file encoding
type Warp struct {
	w io.Writer
}

//AddFile adds a new filw inside warp file
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

	//we need to reset the file to 0, 0, so copy io can performe.
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

//NewWriter creates a new Warp
func NewWriter(w io.Writer) *Warp {
	return &Warp{w: w}
}
