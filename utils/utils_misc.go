package utils

import (
	"bytes"
	"encoding/json"
)

// Misc functions about json operation
// Pretty print of json, debugging purpose
func PrettyString(str string) (string, error) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(str), "", "    "); err != nil {
		return "", err
	}
	return prettyJSON.String(), nil
}

// Find index of string inside array of string
func FindStringIndex(strList []string, target string) int {
	for i, str := range strList {
		if str == target {
			return i
		}
	}
	return -1 // target string not found in the list
}

// Convert array of interfaces into array of strings
func ConvertInterfacesToStrings(interfaces []interface{}) []string {
	strings := make([]string, 0)
	for _, keyInterface := range interfaces {
		strings = append(strings, keyInterface.(string))
	}

	return strings
}

// Compare two string arrays and return leftOnly and rightOnly element
func CompareLists(left []string, right []string) ([]string, []string) {
	leftOnly := []string{}
	rightOnly := []string{}

	for _, leftItem := range left {
		found := false
		for _, rightItem := range right {
			if leftItem == rightItem {
				found = true
				break
			}
		}
		if !found {
			leftOnly = append(leftOnly, leftItem)
		}
	}

	for _, rightItem := range right {
		found := false
		for _, leftItem := range left {
			if rightItem == leftItem {
				found = true
				break
			}
		}
		if !found {
			rightOnly = append(rightOnly, rightItem)
		}
	}

	return leftOnly, rightOnly
}
