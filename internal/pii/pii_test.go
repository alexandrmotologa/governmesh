package pii

import (
	"testing"
)

func TestRegexPatterns(t *testing.T) {
	testText := `
	Here is an email: test.user@enterprise.com
	And a secret: sk-1234567890abcdef1234567890abcdef12
	And an SSN: 123-45-6789
	`
	findings := ScanPromptText(testText, 1)
	if len(findings) < 3 {
		t.Fatalf("expected at least 3 findings, got %d", len(findings))
	}

	foundSSN := false
	for _, f := range findings {
		if f.Type == "ssn" {
			foundSSN = true
			if f.Snippet != "12*******89" {
				t.Errorf("expected redacted SSN, got %s", f.Snippet)
			}
		}
	}
	if !foundSSN {
		t.Errorf("expected to find SSN")
	}
}

func TestVariableRules(t *testing.T) {
	vars := []string{"candidate_ssn", "patient_medical_history", "user_query"}
	findings := ScanVariables(vars, 10)
	if len(findings) != 2 {
		t.Fatalf("expected 2 sensitive variable findings, got %d", len(findings))
	}
}
