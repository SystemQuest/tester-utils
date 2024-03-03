package tester_utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/codecrafters-io/tester-utils/test_case_harness"
	"github.com/stretchr/testify/assert"
)

func passFunc(harness *test_case_harness.TestCaseHarness) error {
	return nil
}

func failFunc(harness *test_case_harness.TestCaseHarness) error {
	return errors.New("fail")
}

func buildTestCasesJson(slugs []string) string {
	testCases := []map[string]string{}

	for index, slug := range slugs {
		testCases = append(testCases, map[string]string{
			"slug":              slug,
			"tester_log_prefix": fmt.Sprintf("test-%d", index+1),
			"title":             fmt.Sprintf("Stage #%d: %s", index+1, slug),
		})
	}

	testCasesJson, _ := json.Marshal(testCases)
	return string(testCasesJson)
}

func buildTester(t *testing.T, definition TesterDefinition, testCasesJson string) Tester {
	tester, err := NewTester(map[string]string{
		"CODECRAFTERS_SUBMISSION_DIR":  "./test_helpers/valid_app_dir",
		"CODECRAFTERS_TEST_CASES_JSON": buildTestCasesJson([]string{"test-1", "test-2"}),
	}, definition)

	if err != nil {
		t.Error(err)
	}

	return tester
}

func TestAllStagesPass(t *testing.T) {
	definition := TesterDefinition{
		TestCases: []TestCase{
			{Slug: "test-1", TestFunc: passFunc},
			{Slug: "test-2", TestFunc: passFunc},
		},
	}

	tester := buildTester(t, definition, buildTestCasesJson([]string{"test-1", "test-2"}))
	exitCode := tester.RunCLI()
	assert.Equal(t, exitCode, 0)
}

func TestOneStageFails(t *testing.T) {
	definition := TesterDefinition{
		TestCases: []TestCase{
			{Slug: "test-1", TestFunc: passFunc},
			{Slug: "test-2", TestFunc: failFunc},
		},
	}

	tester := buildTester(t, definition, buildTestCasesJson([]string{"test-1", "test-2"}))
	exitCode := tester.RunCLI()
	assert.Equal(t, exitCode, 1)
}
