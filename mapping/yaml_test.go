package mapping

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"testing"
)

/**
These patterns will not pass validation.
 */
var badMappings = []struct{
	path string
	redirect string
}{
	{
		"", // empty path
		"https://127.0.0.1",
	},
	{
		"pathA",
		"://127.0.0.1",  // no scheme
	},
	{
		"pathA",  // path has no slash prefix
		"https://127.0.0.1",
	},
	{
		"/pathA",
		"http://127.0.0.1",  // we only accept https, sorry
	},
	{
		"/pathA",
		"ftp://127.0.0.1",  // we only accept https, sorry
	},
}

func Test_MappingValidate(t *testing.T) {
	path := "/"
	redirect := "https://127.0.0.1"
	mapping := Mapping{
		path: redirect,
	}

	if err := mapping.Validate(); err != nil {
		t.Errorf("Could not parse and validate new MappingsFile, error:[%s]", err)
	}
}

func Test_MappingScheme(t *testing.T) {
	path := "/"
	redirect := "https://127.0.0.1"
	mapping := Mapping{
		path: redirect,
	}

	if err := mapping.Validate(); err != nil {
		t.Errorf("Could not parse and validate new MappingsFile, error:[%s]", err)
	}
}

func Test_badMappings(t *testing.T) {
	for index, testData := range badMappings {
		mapping := Mapping{
			testData.path: testData.redirect,
		}
		if err := mapping.Validate(); err == nil {
			msg := fmt.Sprintf("Expected badMappings[%d] to be invalid, ended up being valid.", index)
			t.Errorf(msg)
		}
	}
}

/**
Here we test access to the mappings map. We also enforce that it is a map if anyone changes it.
 */
func Test_MappingsMap(t *testing.T) {
	expectedKey := "test"

	redirectMap := MappingsFile{
		Mappings: map[string]Mapping{
			expectedKey: {
				"/mypath": "https://127.0.0.1",
				"/mypath2": "https://127.0.0.1",
			},
		},
	}

	// GetRedirectUri something we know exists
	if value := redirectMap.GetRedirectUri(expectedKey, "/mypath"); value == "" {
		t.Errorf("Expected a mapping")
	}

	// GetRedirectUri a key that does not exist
	if value := redirectMap.GetRedirectUri("n/a", ""); value != "" {
		t.Errorf("Expected to get an error for a search of key[%s]", "n/a")
	}
}

func Test_MappingFileWithLocalhost(t *testing.T) {
	data := MappingsFile{}
	testFile := `---
mapping:
  localhost:
    "/my-path": https://localhost:8081
    "/": https://localhost:8082
`

	if err := yaml.Unmarshal([]byte(testFile), &data); err != nil {
		t.Errorf("Could not parse test data: %v", err)
	}

	if err := data.Validate(); err == nil {
		t.Errorf("Data was expected to be invalid as you cannot use localhost: %v", err)
	}
}

func Test_MappingFileWithRoot(t *testing.T) {
	data := MappingsFile{}
	testFile := `---
mapping:
  testhost:
    "/my-path": https://localhost:8081
    "/": https://localhost:8082
`

	if err := yaml.Unmarshal([]byte(testFile), &data); err != nil {
		t.Errorf("Could not parse test data: %v", err)
	}

	if err := data.Validate(); err != nil {
		t.Errorf("Data was expected to be valid: %v", err)
	}

	if uri := data.GetRedirectUri("testhost", "/my-path"); uri != "https://localhost:8081" {
		t.Error("Incorrect URI obtained, expected https://localhost:8081")
	}

	if uri := data.GetRedirectUri("testhost", "/"); uri != "https://localhost:8082" {
		t.Error("Incorrect URI obtained, expected https://localhost:8082")
	}

	// we treat root as a wildcard pattern
	if uri := data.GetRedirectUri("testhost", "/something-not-there"); uri != "https://localhost:8082" {
		t.Error("Incorrect URI obtained, expected https://localhost:8082")
	}

}

func Test_MappingFileWithoutRoot(t *testing.T) {
	data := MappingsFile{}
	testFile := `---
mapping:
  testhost:
    "/my-path": https://localhost:8081
`

	if err := yaml.Unmarshal([]byte(testFile), &data); err != nil {
		t.Errorf("Could not parse test data: %v", err)
	}

	if err := data.Validate(); err != nil {
		t.Errorf("Data was expected to be valid: %v", err)
	}

	if uri := data.GetRedirectUri("testhost", "/my-path"); uri != "https://localhost:8081" {
		t.Error("Incorrect URI obtained, expected https://localhost:8081")
	}

	if uri := data.GetRedirectUri("testhost", "/"); uri != "" {
		t.Error("Incorrect URI obtained, expected empty string since mapping doesn't specify a wildcard root '/'")
	}
}