package util

import (
	"bytes"
	"encoding/json"
	"log"
)

func SliceMap[T any, U any](slice []T, function func(T, int) U) []U {
	result := make([]U, len(slice))

	for index, v := range slice {
		result[index] = function(v, index)
	}
	return result
}

func SliceFilter[T any](slice []T, function func(T, int) bool) []T {
	result := make([]T, 0)

	for index, v := range slice {
		if function(v, index) {
			result = append(result, v)
		}
	}
	return result
}

func SliceContains[T comparable](slice []T, value T) bool {
	contains := false
	for _, v := range slice {
		if v == value {
			contains = true
			break
		}
	}
	return contains
}

func ConvertMapToJsonString(m map[string]any) (string, error) {
	jsonBytes, err := json.Marshal(m)
	if err != nil {
		log.Printf("Error encoding in JSON: %s", err)
		return "", err
	}

	var buffer bytes.Buffer
	err = json.Indent(&buffer, jsonBytes, "", "\t")
	if err != nil {
		log.Printf("Error pretty printing in JSON: %s", err)
		return "", err
	}

	return buffer.String(), nil
}
