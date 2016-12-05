package warpdrive

import (
	"fmt"
	"io/ioutil"
	"os"
)

func loadFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Printf("%s", string(b))
}

func saveFile(filename, content string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		fmt.Println(err.Error())
	}
}
