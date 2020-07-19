package main

import (
	"reflect"
	"testing"
)

func TestGetDiseaseListFromFile(t *testing.T) {
	tests := map[string]struct {
		filename  string
		shouldErr bool
	}{
		"valid filename with valid JSON":   {"test-fixtures/diseases.json", false},
		"valid filename with invalid JSON": {"test-fixtures/malformed.json", true},
		"invalid filename":                 {"test-fixtures/foobarbaz", true},
	}

	testDiseases := new(diseaseList)
	testDiseases.Diseases = []string{"foo", "bar"}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result, err := getDiseaseListFromFile(test.filename)
			if test.shouldErr && err == nil {
				t.Fatalf("Expected error, but no error occured")
			}

			if !test.shouldErr && err != nil {
				t.Fatalf("Expected no error, but got: %s", err.Error())
			}

			if !reflect.DeepEqual(result, testDiseases) && !test.shouldErr {
				t.Fatalf("Results did not match\nGot: %+v\nExpected: %+v", result, testDiseases)
			}
		})
	}
}

func TestGetDisease(t *testing.T) {
	testDiseases := new(diseaseList)
	testDiseases.Diseases = []string{"foo", "bar", "baz"}

	testMalformed := new(diseaseList)
	testMalformed.Diseases = []string{}

	tests := map[string]struct {
		input     *diseaseList
		expected  []string
		shouldErr bool
	}{
		"valid diseaseList":   {testDiseases, []string{"foo", "bar", "baz"}, false},
		"invalid diseaseList": {testMalformed, []string{}, true},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result, err := test.input.getDisease()
			if test.shouldErr && err == nil {
				t.Fatalf("Expected error, but no error occurred")
			}

			if !test.shouldErr && err != nil {
				t.Fatalf("Expected no error, but got: %s", err.Error())
			}

			if !test.shouldErr && !stringInSlice(result, test.expected) {
				t.Fatalf("Result not in expect list\nGot: %s\nExpected one from: %v", result, test.expected)
			}
		})
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if a == b {
			return true
		}
	}

	return false
}
