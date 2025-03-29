package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

func TestGetS3File(t *testing.T) {
	s3, _ := InitialiseS3()
	objs, err := getJsonFile(s3, "jsonfiledatacache", "test.json")
	t.Log(objs)
	t.Log(err)
}

func TestBye(t *testing.T) {
	byteArray := []byte(`{"name": "Alice", "age": 30, "active": true}`)
	fmt.Println(byteArray)
	var result map[string]any
	decoder := json.NewDecoder(bytes.NewReader(byteArray))

	// Decode JSON efficiently
	if err := decoder.Decode(&result); err != nil {
		log.Fatal("❌ JSON Decoding Error:", err)
	}

	// Print the resulting map
	fmt.Println("✅ Decoded JSON:", result)

}

type Person struct {
	Name   string `json:"name"`
	Age    int    `json:"age"`
	ByeBye *struct {
		Test int `json:"test"`
	} `json:"byebye"`
}

func TestJSONValid(t *testing.T) {
	// Example JSON with some invalid objects
	jsonString := `[{"name": "Alice", "age": 30}, {"name": "Bob"}, {"name": "Charlie", "age": "old"}, {"name": "David", "age": 40}]`

	// Convert JSON string to a byte reader
	reader := bytes.NewReader([]byte(jsonString))

	// Use a JSON decoder for efficient streaming
	decoder := json.NewDecoder(reader)

	// Expecting a JSON array (must start with `[`)
	_, err := decoder.Token() // Read opening '['
	if err != nil {
		log.Fatal("❌ JSON Decoding Error:", err)
	}

	var validPeople []Person // Store valid struct objects

	// Loop through each object in the JSON array
	for decoder.More() {
		var p Person
		err := decoder.Decode(&p) // Try to decode into the struct
		if err == nil {
			validPeople = append(validPeople, p) // Only add valid objects
		} else {
			fmt.Println("⚠️ Skipping invalid object:", err)
		}
	}

	fmt.Println("✅ Decoded Structs:", validPeople)
}
