package warp

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Compress gets a map of <Target>:Source values. Target is the actual name and src is
// hash name path. So we need to load the src first and change the name of the file to target
func Compress(paths map[string]string, output io.Writer) error {
	fileWriter := gzip.NewWriter(output)
	defer fileWriter.Close()

	tarfileWriter := tar.NewWriter(fileWriter)
	defer tarfileWriter.Close()

	for src, target := range paths {
		file, err := os.Open(src)
		if err != nil {
			return err
		}

		// we need to refactor this code,
		// look into http://stackoverflow.com/questions/24720097/golang-defer-behavior for more info
		defer file.Close()

		fileInfo, err := file.Stat()
		if err != nil {
			return err
		}

		if fileInfo.IsDir() {
			continue
		}

		header := new(tar.Header)
		header.Name = target
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
	}

	return nil
}

// Extract extracts input stream into targetPath
func Extract(input io.ReadCloser, targetPath string) error {
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
		// and join it with targetPath
		filename := filepath.Join(targetPath, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			// handle directory
			err = os.MkdirAll(filename, os.ModePerm)
			if err != nil {
				return err
			}

		case tar.TypeReg:
		case tar.TypeRegA:
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
				return err
			}

			writer.Close()

		default:
			return fmt.Errorf("Unable to untar type : %c in file %s", header.Typeflag, filename)
		}
	}

	return nil
}
