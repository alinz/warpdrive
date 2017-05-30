package helper

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// GetAllFiles returns the path to all files, only files not directories, inside a folder
func GetAllFiles(target string) ([]string, error) {
	fileList := []string{}

	// the following tw lines adds the "/" to the end of path
	// this is helpful to create relative path
	target = filepath.Join(target, "/")
	target = target + "/"

	err := filepath.Walk(target, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			fileList = append(fileList, strings.Replace(path, target, "", -1))
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return fileList, nil
}

// BundleCompress is getting a path to bundle folder write the
// compress tar.gz to the given output
func BundleCompress(rootPath string, output io.Writer) error {
	fileWriter := gzip.NewWriter(output)
	defer fileWriter.Close()

	tarfileWriter := tar.NewWriter(fileWriter)
	defer tarfileWriter.Close()

	files, err := GetAllFiles(rootPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		// in order to prevent defer stacking up,
		// a function wrapper is being used. so once the
		// file is close, defer will be executed to close the
		// file handler
		err := func(path string) error {
			file, err := os.Open(filepath.Join(rootPath, path))
			if err != nil {
				return err
			}

			defer file.Close()

			fileInfo, err := file.Stat()
			if err != nil {
				return err
			}

			header := new(tar.Header)
			header.Name = path
			header.Size = fileInfo.Size()
			header.Mode = int64(fileInfo.Mode())
			header.ModTime = fileInfo.ModTime()

			err = tarfileWriter.WriteHeader(header)
			if err != nil {
				return err
			}

			_, err = io.Copy(tarfileWriter, file)
			if err != nil {
				return err
			}

			return nil
		}(file)

		if err != nil {
			return err
		}
	}

	return nil
}

// BundleUncompress reads the stream of tar.gz file and decompress it to given path
func BundleUncompress(input io.Reader, path string) error {
	fileReader, err := gzip.NewReader(input)
	if err != nil {
		return err
	}

	defer fileReader.Close()

	tarBallReader := tar.NewReader(fileReader)

	for {
		header, err := tarBallReader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		// get the individual filename and extract to the current directory
		// and join it with path
		filename := filepath.Join(path, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			// handle directory
			err = os.MkdirAll(filename, os.ModePerm)
			if err != nil {
				return err
			}

		case tar.TypeReg, tar.TypeRegA:
			// we need to make sure the folder exists
			dir := filepath.Dir(filename)
			err = os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				return err
			}

			// handle normal file
			writer, err := os.Create(filename)
			if err != nil {
				return err
			}

			io.Copy(writer, tarBallReader)

			err = os.Chmod(filename, os.ModePerm)
			if err != nil {
				writer.Close()
				return err
			}

			writer.Close()

		default:
			return fmt.Errorf("Unable to untar type : %c in file %s", header.Typeflag, filename)
		}
	}

	return nil
}
