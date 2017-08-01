// TestSuite.go
package autositespeed

import (
	"encoding/json"
	"io/ioutil"
)

type UrlTest struct {
	Name      string `json:"name"`
	Url       string `json:"url"`
	PreScript string `json:"preScript"`
}

type TestSuite struct {
	SuiteName string    `json:"suiteName"`
	Tests     []UrlTest `json:"tests"`
}

func ReadTestSuiteFromFile(testSuiteFile string) (*TestSuite, error) {
	fileBytes, err := ioutil.ReadFile(testSuiteFile)
	var testSuite TestSuite

	if err == nil {
		err = json.Unmarshal(fileBytes, &testSuite)
	}

	return &testSuite, err
}
