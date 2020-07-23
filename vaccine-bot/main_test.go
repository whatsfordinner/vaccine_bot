package main

import (
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

func TestGetParameter(t *testing.T) {
	tests := map[string]struct {
		paramSetInEnv bool
		paramSetInSSM bool
		shouldErr     bool
	}{
		"parameter set in environment": {true, false, false},
		"parameter set in SSM":         {false, true, true},
		"parameter not set":            {false, false, true},
	}

	expected := "foobar"

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if test.paramSetInEnv {
				os.Setenv("TEST_PARAM", "foobar")
			} else {
				os.Unsetenv("TEST_PARAM")
			}

			if test.paramSetInSSM {
				// not yet implemented
			}

			result, err := getParameter("TEST_PARAM")

			if test.shouldErr && err == nil {
				t.Fatalf("Expecter error, but no error occurred")
			}

			if !test.shouldErr && err != nil {
				t.Fatalf("Expected no error, but got: %s", err.Error())
			}

			if result != expected && err == nil {
				t.Fatalf("Results did not match\nGot: %s\nExpected: %s", result, expected)
			}
		})
	}
}

func TestGetTwitterConfig(t *testing.T) {
	tests := map[string]struct {
		consumerKeySet    bool
		consumerSecretSet bool
		accessTokenSet    bool
		accessSecretSet   bool
		shouldErr         bool
	}{
		"all variables set":       {true, true, true, true, false},
		"consumer key not set":    {false, true, true, true, true},
		"consumer secret not set": {true, false, true, true, true},
		"access token not set":    {true, true, false, true, true},
		"access secret not set":   {true, true, true, false, true},
	}

	expected := &twitterConfig{
		consumerKey:    "foo",
		consumerSecret: "foo",
		accessToken:    "foo",
		accessSecret:   "foo",
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if test.consumerKeySet {
				os.Setenv("TWITTER_CONSUMER_KEY", "foo")
			} else {
				os.Unsetenv("TWITTER_CONSUMER_KEY")
			}

			if test.consumerSecretSet {
				os.Setenv("TWITTER_CONSUMER_SECRET", "foo")
			} else {
				os.Unsetenv("TWITTER_CONSUMER_SECRET")
			}

			if test.accessTokenSet {
				os.Setenv("TWITTER_ACCESS_TOKEN", "foo")
			} else {
				os.Unsetenv("TWITTER_ACCESS_TOKEN")
			}

			if test.accessSecretSet {
				os.Setenv("TWITTER_ACCESS_SECRET", "foo")
			} else {
				os.Unsetenv("TWITTER_ACCESS_SECRET")
			}

			result, err := getTwitterConfig()

			if test.shouldErr && err == nil {
				t.Fatalf("Expected error, but no error occurred")
			}

			if !test.shouldErr && err != nil {
				t.Fatalf("Expected no error, but got: %s", err.Error())
			}

			if !reflect.DeepEqual(result, expected) && !test.shouldErr {
				t.Fatalf("Results did not match\nGot: %+v\nExpected: %+v", result, expected)
			}
		})
	}
}

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
				t.Fatalf("Result not in expected list\nGot: %s\nExpected one from: %v", result, test.expected)
			}
		})
	}
}

func TestBuildTweet(t *testing.T) {
	testDisease := "foobar"
	expected := "Vaccinate your kids against foobar"
	result := buildTweet(testDisease)

	if result != expected {
		t.Fatalf("Result did not match expected\nGot: %s\nExpected: %s", result, expected)
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
