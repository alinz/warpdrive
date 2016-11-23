package web

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
)

// WriteAsJSON a helper function to simplifies the JSON serilization
func WriteAsJSON(w http.ResponseWriter, response interface{}, status int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	if response != nil {
		json.NewEncoder(w).Encode(response)
	}
}

// WriteAsText a helper function to simplifies the string serilization
func WriteAsText(w http.ResponseWriter, response interface{}, status int) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(status)
	w.Write([]byte(fmt.Sprintf("%v", response)))
}

// StreamJSONToStruct converts stream of json to a defined struct
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

// StreamJSONToStructWithLimit similar to StreamJSONToStruct but with limit in
// payload size
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

// JSONValidation accept any struct with json tag. it will returns an error
// if any field inside that struct has `required` tag value and becomes nil
func JSONValidation(jsonObj interface{}) error {
	tagType := reflect.TypeOf(jsonObj)
	jsonFieldName := ""

	// if the interface is pointer we need to get access to actual value
	interfaceIsPoniter := tagType.Kind() == reflect.Ptr
	if interfaceIsPoniter {
		tagType = tagType.Elem()
	}

	// loops through all field
	for i := 0; i < tagType.NumField(); i++ {
		field := tagType.Field(i)
		jsonTag := field.Tag.Get("json")
		jsonTagValues := strings.Split(jsonTag, ",")

		// the first arguments in jsonFildName is always json represenation of field
		if len(jsonTagValues) > 0 {
			jsonFieldName = jsonTagValues[0]
		}

		// we are searching inside json tag's value to see if `required` is presented.
		isRequired := false
		for _, jsonTagValue := range jsonTagValues {
			if jsonTagValue == "required" {
				isRequired = true
				break
			}
		}

		// if require is presented, we check whether the value of that field is
		// nil or not. Remember, in order for this function to work, all fields in struct
		// needs to be converted into poniter instaed of value.
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

// BodyParser loads builder with maxSize and tries to load the message.
// if for some reason it can't parse the message, it will return an error.
// if successful, it will put the processed data into context with key 'json_body'
func BodyParser(body interface{}, maxSize int64) func(next http.Handler) http.Handler {
	// body needs to be a pointer type
	bodyType := reflect.TypeOf(body).Elem()
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// to is a new instance of type of body. It is an interface.
			to := reflect.New(bodyType).Interface()

			if err := StreamJSONToStructWithLimit(r.Body, to, maxSize); err != nil {
				Respond(w, http.StatusUnprocessableEntity, err)
				return
			}

			// check for required fields
			if err := JSONValidation(to); err != nil {
				Respond(w, http.StatusBadRequest, err)
				return
			}

			ctx := context.WithValue(r.Context(), "parsed:body", to)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
