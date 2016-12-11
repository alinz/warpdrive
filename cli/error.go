package cli

import (
	"encoding/json"
	"fmt"
	"io"
)

func parseErrorMessage(body io.Reader) error {
	var err error

	errorMessage := struct {
		Error string `json:"error"`
	}{}

	err = json.NewDecoder(body).Decode(&errorMessage)
	if err != nil {
		errorMessage.Error = err.Error()
	}

	return fmt.Errorf(errorMessage.Error)
}
