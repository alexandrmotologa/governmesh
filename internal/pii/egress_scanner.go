package pii

import (
	"fmt"
	"strings"
)

// ScanPromptText scans a raw string (e.g. system prompt, template content) for literal PII patterns.
func ScanPromptText(text string, baseLine int) []Finding {
	var findings []Finding
	lines := strings.Split(text, "\n")

	for i, line := range lines {
		currentLine := baseLine + i
		for _, pattern := range RegexPatterns {
			matches := pattern.Regex.FindAllStringIndex(line, -1)
			for _, m := range matches {
				rawMatch := line[m[0]:m[1]]
				redacted := redactSnippet(rawMatch)

				findings = append(findings, Finding{
					Type:        pattern.Type,
					Snippet:     redacted,
					Line:        currentLine,
					Severity:    pattern.Severity,
					Description: pattern.Description,
				})
			}
		}
	}

	return findings
}

// ScanVariables checks extracted variable placeholders (e.g. from prompt templates or format calls).
func ScanVariables(vars []string, line int) []Finding {
	var findings []Finding

	for _, v := range vars {
		if rule := MatchVariable(v); rule != nil {
			findings = append(findings, Finding{
				Type:        rule.Type,
				Snippet:     fmt.Sprintf("{%s}", v),
				Line:        line,
				Severity:    rule.Severity,
				Variable:    v,
				Description: rule.Description,
			})
		}
	}

	return findings
}

// DetermineEgressTarget identifies whether a model provider is external cloud, local runtime, or VPC.
func DetermineEgressTarget(provider, endpoint string) string {
	lowerEndpoint := strings.ToLower(endpoint)
	lowerProvider := strings.ToLower(provider)

	if strings.Contains(lowerEndpoint, "localhost") ||
		strings.Contains(lowerEndpoint, "127.0.0.1") ||
		strings.Contains(lowerEndpoint, "ollama") ||
		strings.Contains(lowerEndpoint, "vllm") ||
		strings.Contains(lowerProvider, "local") ||
		strings.Contains(lowerProvider, "ollama") {
		return "local_runtime"
	}

	if strings.Contains(lowerEndpoint, "internal.") ||
		strings.Contains(lowerEndpoint, ".corp") ||
		strings.Contains(lowerEndpoint, ".local") {
		return "private_vpc"
	}

	return "external_cloud"
}

// redactSnippet masks characters to protect sensitive data while preserving format for auditors.
func redactSnippet(s string) string {
	runes := []rune(s)
	length := len(runes)
	if length <= 4 {
		return "****"
	}
	prefix := string(runes[:2])
	suffix := string(runes[length-2:])
	stars := strings.Repeat("*", length-4)
	return prefix + stars + suffix
}
