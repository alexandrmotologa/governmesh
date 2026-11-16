package pii

import (
	"regexp"
	"strings"
)

// PatternDefinition specifies a sensitive data pattern match.
type PatternDefinition struct {
	Type        string
	Regex       *regexp.Regexp
	Severity    string
	Description string
}

// VariableRule defines keyword identifiers in variables that signify sensitive inputs.
type VariableRule struct {
	Type        string
	Keywords    []string
	Severity    string
	Description string
}

var (
	// Pre-compiled regex patterns for sensitive literals.
	RegexPatterns = []PatternDefinition{
		{
			Type:        "email",
			Regex:       regexp.MustCompile(`(?i)\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}\b`),
			Severity:    "medium",
			Description: "Email address detected in prompt context",
		},
		{
			Type:        "ssn",
			Regex:       regexp.MustCompile(`\b\d{3}-\d{2}-\d{4}\b`),
			Severity:    "critical",
			Description: "US Social Security Number pattern detected",
		},
		{
			Type:        "credit_card",
			Regex:       regexp.MustCompile(`\b(?:4[0-9]{12}(?:[0-9]{3})?|5[1-5][0-9]{14}|3[47][0-9]{13}|6(?:011|5[0-9]{2})[0-9]{12})\b`),
			Severity:    "critical",
			Description: "Payment card number detected",
		},
		{
			Type:        "iban",
			Regex:       regexp.MustCompile(`\b[A-Z]{2}\d{2}[A-Z0-9]{4}\d{7}([A-Z0-9]?){0,16}\b`),
			Severity:    "critical",
			Description: "International Bank Account Number (IBAN) detected",
		},
		{
			Type:        "api_key_openai",
			Regex:       regexp.MustCompile(`sk-[a-zA-Z0-9]{32,}`),
			Severity:    "critical",
			Description: "Hardcoded OpenAI API Secret Key detected",
		},
		{
			Type:        "api_key_anthropic",
			Regex:       regexp.MustCompile(`sk-ant-[a-zA-Z0-9]{32,}`),
			Severity:    "critical",
			Description: "Hardcoded Anthropic API Secret Key detected",
		},
	}

	// VariableRules inspects prompt variable placeholders like {patient_medical_history} or ${candidate_ssn}.
	VariableRules = []VariableRule{
		{
			Type:        "phi_medical",
			Keywords:    []string{"patient", "diagnosis", "prescription", "medical_history", "clinical_record", "health_condition", "treatment_plan", "symptoms"},
			Severity:    "critical",
			Description: "Protected Health Information (PHI) variable binding",
		},
		{
			Type:        "pii_national_id",
			Keywords:    []string{"ssn", "social_security", "national_id", "passport_number", "id_card", "cnp", "tax_id"},
			Severity:    "critical",
			Description: "National identification number variable binding",
		},
		{
			Type:        "financial",
			Keywords:    []string{"salary", "income", "credit_score", "bank_account", "card_number", "routing_number", "iban", "net_worth", "debt"},
			Severity:    "high",
			Description: "Financial data or credit rating variable binding",
		},
		{
			Type:        "biometric",
			Keywords:    []string{"biometric", "facial_embedding", "face_scan", "fingerprint", "retina_scan", "voiceprint"},
			Severity:    "critical",
			Description: "Biometric identifier variable binding",
		},
		{
			Type:        "sensitive_demographic",
			Keywords:    []string{"sexual_orientation", "religion", "ethnic_origin", "political_affiliation", "trade_union"},
			Severity:    "critical",
			Description: "Special category data (EU AI Act Article 5 / GDPR Article 9)",
		},
		{
			Type:        "pii_personal",
			Keywords:    []string{"user_email", "phone_number", "home_address", "date_of_birth", "dob", "full_name"},
			Severity:    "medium",
			Description: "Personal contact and identity information binding",
		},
	}
)

// MatchVariable checks if a variable name matches any sensitive rule keywords.
func MatchVariable(varName string) *VariableRule {
	normalized := strings.ToLower(varName)
	// Remove brackets or punctuation
	normalized = strings.Trim(normalized, "{}$_:- ")

	for _, rule := range VariableRules {
		for _, kw := range rule.Keywords {
			if strings.Contains(normalized, kw) {
				return &rule
			}
		}
	}
	return nil
}
