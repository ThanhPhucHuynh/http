package http

import (
	"testing"
)

func TestToJson(t *testing.T) {
	type TestData struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	// Test case 1: Valid input
	data := TestData{Name: "John", Age: 30}
	expectedJSON := `{"name":"John","age":30}`
	jsonData, err := ToJson(data)
	if err != nil {
		t.Errorf("ToJson returned an error for valid input: %v", err)
	}
	if jsonData != expectedJSON {
		t.Errorf("ToJson result is incorrect. Expected: %s, Got: %s", expectedJSON, jsonData)
	}

	// Test case 2: Error case
	invalidData := make(chan int) // This type cannot be marshaled to JSON
	_, err = ToJson(invalidData)
	if err == nil {
		t.Error("ToJson should return an error for invalid input")
	}
}

func TestToReader(t *testing.T) {
	// Test case 1: Non-empty bodyString
	bodyString := "Hello, world!"
	reader := ToReader(bodyString)
	buffer := make([]byte, len(bodyString))
	n, err := reader.Read(buffer)
	if err != nil {
		t.Errorf("ToReader returned an error: %v", err)
	}
	if n != len(bodyString) {
		t.Errorf("ToReader did not read the correct number of bytes. Expected: %d, Got: %d", len(bodyString), n)
	}
	if string(buffer) != bodyString {
		t.Errorf("ToReader result is incorrect. Expected: %s, Got: %s", bodyString, string(buffer))
	}

	// Test case 2: Empty bodyString
	emptyReader := ToReader("")
	emptyBuffer := make([]byte, 1)
	n, err = emptyReader.Read(emptyBuffer)
	if err == nil {
		t.Error("ToReader should return an error for an empty bodyString")
	}
}
