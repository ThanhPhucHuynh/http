package http

import (
	"bytes"
	"encoding/json"
)

const EMPTY_STRING = ""

func ToJson(data interface{}) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return EMPTY_STRING, err
	}
	return string(jsonData), nil
}

func ToReader(bodyString string) *bytes.Reader {
	return bytes.NewReader([]byte(bodyString))
}
