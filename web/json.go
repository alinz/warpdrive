package web

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
)

//WriteAsJSON a helper function to simplifies the JSON serilization
func WriteAsJSON(w http.ResponseWriter, response interface{}, status int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	if response != nil {
		json.NewEncoder(w).Encode(response)
	}
}

//StreamJSONToStruct converts stream of json to a defined struct
func StreamJSONToStruct(r io.Reader, v interface{}) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	return nil
}

//StreamJSONToStructWithLimit similar to StreamJSONToStruct but with limit in
//payload size
func StreamJSONToStructWithLimit(r io.Reader, v interface{}, limit int64) error {
	raw, err := ioutil.ReadAll(io.LimitReader(r, limit))

	if err != nil {
		return err
	}

	if err := json.Unmarshal(raw, &v); err != nil {
		return err
	}

	return nil
}

//JSONValidation accept any struct with json tag. it will returns an error
//if any field inside that struct has `required` tag value and becomes nil
func JSONValidation(jsonObj interface{}) error {
	tagType := reflect.TypeOf(jsonObj)
	jsonFieldName := ""

	//if the interface is pointer we need to get access to actual value
	interfaceIsPoniter := tagType.Kind() == reflect.Ptr
	if interfaceIsPoniter {
		tagType = tagType.Elem()
	}

	//loops through all field
	for i := 0; i < tagType.NumField(); i++ {
		field := tagType.Field(i)
		jsonTag := field.Tag.Get("json")
		jsonTagValues := strings.Split(jsonTag, ",")

		//the first arguments in jsonFildName is always json represenation of field
		if len(jsonTagValues) > 0 {
			jsonFieldName = jsonTagValues[0]
		}

		//we are searching inside json tag's value to see if `required` is presented.
		isRequired := false
		for _, jsonTagValue := range jsonTagValues {
			if jsonTagValue == "required" {
				isRequired = true
				break
			}
		}

		//if require is presented, we check whether the value of that field is
		//nil or not. Remember, in order for this function to work, all fields in struct
		//needs to be converted into poniter instaed of value.
		if isRequired {
			immutable := reflect.ValueOf(jsonObj)
			//if the interface is pointer we need to get access to actual value
			if interfaceIsPoniter {
				immutable = immutable.Elem()
			}
			if immutable.FieldByName(field.Name).IsNil() {
				return fmt.Errorf("field '%s' required", jsonFieldName)
			}
		}
	}

	return nil
}
