package godogs

import (
	"context"
	"errors"
	"fmt"
	"testing"
  "os"
	"github.com/cucumber/godog"
)

// godogsCtxKey is the key used to store the available godogs in the context.Context.
type godogsCtxKey struct{}

func thereAreGodogs(ctx context.Context, available int) (context.Context, error) {
	return context.WithValue(ctx, godogsCtxKey{}, available), nil
}

func iEat(ctx context.Context, num int) (context.Context, error) {
	available, ok := ctx.Value(godogsCtxKey{}).(int)
	if !ok {
		return ctx, errors.New("there are no godogs available")
	}

	if available < num {
		return ctx, fmt.Errorf("you cannot eat %d godogs, there are %d available", num, available)
	}

	available -= num

	return context.WithValue(ctx, godogsCtxKey{}, available), nil
}

func thereShouldBeRemaining(ctx context.Context, remaining int) error {
	available, ok := ctx.Value(godogsCtxKey{}).(int)
	if !ok {
		return errors.New("there are no godogs available")
	}

	if available != remaining {
		return fmt.Errorf("expected %d godogs to be remaining, but there is %d", remaining, available)
	}

	return nil
}

// TestFeatures is the entry point for running your feature tests.

func TestFeatures(t *testing.T) {
	// Create a JSON file to store the output
	jsonOutputFile, err := os.Create("test/report/cucumber_report.json")
	if err != nil {
		t.Fatalf("could not create json output file: %v", err)
	}
	defer jsonOutputFile.Close()

	// Run the Godog test suite
	status := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,  // Your step definitions
		Options: &godog.Options{
			Format: "cucumber",              // Generate cucumber JSON output
			Paths:  []string{"features"},    // Path to your feature files
			Output: jsonOutputFile,          // Output to JSON file
			TestingT: t,
		},
	}.Run()

	if status != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}


// InitializeScenario initializes your scenario steps.
func InitializeScenario(sc *godog.ScenarioContext) {
	sc.Step(`^there are (\d+) godogs$`, thereAreGodogs)
	sc.Step(`^I eat (\d+)$`, iEat)
	sc.Step(`^there should be (\d+) remaining$`, thereShouldBeRemaining)
}
