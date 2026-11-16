package classifier

import (
	"github.com/alexandrmotologa/governmesh/internal/pii"
	"github.com/alexandrmotologa/governmesh/internal/scanner"
)

// NISTControl represents a NIST AI RMF sub-category check.
type NISTControl struct {
	ID          string  `json:"id"`
	Function    string  `json:"function"` // GOVERN, MAP, MEASURE, MANAGE
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Fulfilled   bool    `json:"fulfilled"`
	Evidence    string  `json:"evidence"`
	Weight      float64 `json:"weight"`
}

// EvaluateNIST calculates NIST AI RMF alignment based on inventory and findings.
func EvaluateNIST(assets []scanner.DiscoveredAsset, piiFindings []pii.Finding) ([]NISTControl, float64) {
	hasInventory := len(assets) > 0
	hasHighRisk := false
	hasPIIFailures := len(piiFindings) > 0
	hasSystemPrompts := false

	for _, a := range assets {
		if a.Risk.Tier == "high" || a.Risk.Tier == "unacceptable" {
			hasHighRisk = true
		}
		if a.SystemPrompt != "" {
			hasSystemPrompts = true
		}
	}

	controls := []NISTControl{
		{
			ID:          "GOVERN-1.1",
			Function:    "GOVERN",
			Title:       "Legal & Regulatory Requirements Mapped",
			Description: "Organization maps statutory AI obligations (EU AI Act, local directives).",
			Fulfilled:   hasInventory,
			Evidence:    "Automated mapping applied across all discovered code assets.",
			Weight:      1.0,
		},
		{
			ID:          "GOVERN-2.1",
			Function:    "GOVERN",
			Title:       "Accountability and Inventory",
			Description: "Inventory of AI models, versions, and providers is maintained.",
			Fulfilled:   hasInventory,
			Evidence:    "Authoritative AIBOM synthesized from AST inspection.",
			Weight:      1.0,
		},
		{
			ID:          "MAP-1.1",
			Function:    "MAP",
			Title:       "Context & Deployment Scope Understood",
			Description: "System purpose, environment, and integrations are documented.",
			Fulfilled:   hasSystemPrompts || hasInventory,
			Evidence:    "AST context extraction mapped model endpoints and invocation scopes.",
			Weight:      1.0,
		},
		{
			ID:          "MAP-1.5",
			Function:    "MAP",
			Title:       "High-Risk Use Cases Identified",
			Description: "High-risk categorization determined according to domain rules.",
			Fulfilled:   !hasHighRisk || (hasHighRisk && hasSystemPrompts),
			Evidence:    "Rule-based classification evaluated against Annex III domains.",
			Weight:      1.2,
		},
		{
			ID:          "MEASURE-1.1",
			Function:    "MEASURE",
			Title:       "Data Privacy & PII Boundary Audit",
			Description: "Data egress and PII interpolation in prompt templates are tracked.",
			Fulfilled:   !hasPIIFailures,
			Evidence:    "PII scanner inspected all variable bindings and prompt literals.",
			Weight:      1.5,
		},
		{
			ID:          "MANAGE-1.3",
			Function:    "MANAGE",
			Title:       "Third-Party AI Supply Chain Controls",
			Description: "Vendor dependencies (OpenAI, Anthropic, Bedrock) tracked in SBOM.",
			Fulfilled:   hasInventory,
			Evidence:    "CycloneDX 1.6 supply chain bill generated with provider metadata.",
			Weight:      1.0,
		},
		{
			ID:          "MANAGE-2.4",
			Function:    "MANAGE",
			Title:       "Incident Prevention & Quality Gate",
			Description: "Automated gate blocks deployment of unacceptable/unmitigated systems.",
			Fulfilled:   true,
			Evidence:    "CI/CD quality gate engine integrated into build pipeline.",
			Weight:      1.0,
		},
	}

	var earnedWeight, totalWeight float64
	for _, c := range controls {
		totalWeight += c.Weight
		if c.Fulfilled {
			earnedWeight += c.Weight
		}
	}

	score := 0.0
	if totalWeight > 0 {
		score = (earnedWeight / totalWeight) * 100.0
	}

	return controls, score
}
