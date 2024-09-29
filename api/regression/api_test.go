package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"os"
	"github.com/cucumber/godog"
)

type apiFeature struct {
	resp *httptest.ResponseRecorder
}

// Reset the response recorder before each scenario
func (a *apiFeature) resetResponse(*godog.Scenario) {
	a.resp = httptest.NewRecorder()
}

// Simulate sending a request to an endpoint (like /version)
func (a *apiFeature) iSendrequestTo(method, endpoint string) (err error) {
	req, err := http.NewRequest(method, endpoint, nil)
	if err != nil {
		return
	}

	// handle panic to ensure tests don’t crash
	defer func() {
		switch t := recover().(type) {
		case string:
			err = fmt.Errorf(t)
		case error:
			err = t
		}
	}()

	switch endpoint {
	case "/version":
		getVersion(a.resp, req)
	default:
		err = fmt.Errorf("unknown endpoint: %s", endpoint)
	}
	return
}

// Check that the response code matches the expected value
func (a *apiFeature) theResponseCodeShouldBe(code int) error {
	if code != a.resp.Code {
		return fmt.Errorf("expected response code to be: %d, but actual is: %d", code, a.resp.Code)
	}
	return nil
}

// Validate that the response body matches the expected JSON
func (a *apiFeature) theResponseShouldMatchJSON(body *godog.DocString) (err error) {
	var expected, actual interface{}

	// Parse the expected JSON
	if err = json.Unmarshal([]byte(body.Content), &expected); err != nil {
		return
	}

	// Parse the actual response JSON
	if err = json.Unmarshal(a.resp.Body.Bytes(), &actual); err != nil {
		return
	}

	// Check if expected and actual JSON match
	if !reflect.DeepEqual(expected, actual) {
		return fmt.Errorf("expected JSON does not match actual, %v vs. %v", expected, actual)
	}
	return nil
}

// Initialize the test suite and define step mappings
func InitializeScenario(ctx *godog.ScenarioContext) {
    api := &apiFeature{}

    ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
        api.resetResponse(sc)
        return ctx, nil
    })

    ctx.Step(`^I send "(GET|POST|PUT|DELETE)" request to "([^"]*)"$`, api.iSendrequestTo)
    ctx.Step(`^the response code should be (\d+)$`, api.theResponseCodeShouldBe)
    ctx.Step(`^the response should match json:$`, api.theResponseShouldMatchJSON)
}


// Run the test suite
func TestAPIFeatureSuite(t *testing.T) {
	// Create the JSON output file for the report
	jsonOutputFile, err := os.Create("test/report/cucumber_report.json")
	if err != nil {
		t.Fatalf("could not create json output file: %v", err)
	}
	defer jsonOutputFile.Close()

	status := godog.TestSuite{
		ScenarioInitializer: InitializeScenario, // Your step definitions
		Options: &godog.Options{
			Format: "cucumber",               // JSON output format
			Paths:  []string{"/Users/tim.sulaimanov/Documents/go_dog/features/godogs.feature"},     // Correct path to feature files
			Output: jsonOutputFile,           // Where the report is stored
			TestingT: t,
		},
	}.Run()
	

	if status != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
