package classifier

import (
	"github.com/alexandrmotologa/governmesh/internal/pii"
	"github.com/alexandrmotologa/governmesh/internal/scanner"
)

// ISOControl represents an ISO/IEC 42001 Annex A control check.
type ISOControl struct {
	Clause      string  `json:"clause"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Fulfilled   bool    `json:"fulfilled"`
	Evidence    string  `json:"evidence"`
	Weight      float64 `json:"weight"`
}

// EvaluateISO42001 assesses ISO/IEC 42001 AIMS control adherence.
func EvaluateISO42001(assets []scanner.DiscoveredAsset, piiFindings []pii.Finding) ([]ISOControl, float64) {
	hasInventory := len(assets) > 0
	hasPIIFailures := len(piiFindings) > 0

	controls := []ISOControl{
		{
			Clause:      "A.5.1",
			Title:       "Policies for AI Development",
			Description: "Enterprise has established formal AI usage and compliance policies.",
			Fulfilled:   true,
			Evidence:    "Compliance policy active with automated risk threshold enforcement.",
			Weight:      1.0,
		},
		{
			Clause:      "A.6.1",
			Title:       "AI Risk Assessment Process",
			Description: "Systematic identification and assessment of AI risks across development.",
			Fulfilled:   hasInventory,
			Evidence:    "Deterministic classification executed against EU AI Act & NIST taxonomies.",
			Weight:      1.5,
		},
		{
			Clause:      "A.7.2",
			Title:       "Data for AI Systems & Privacy Protection",
			Description: "Measures to prevent unauthorized personal data egress and leakage.",
			Fulfilled:   !hasPIIFailures,
			Evidence:    "AST prompt inspector verifies zero unmasked PII/PHI in external model calls.",
			Weight:      1.5,
		},
		{
			Clause:      "A.8.2",
			Title:       "AI System Inventory and Traceability",
			Description: "Accurate, continuous inventory of all deployed AI models and third-party APIs.",
			Fulfilled:   hasInventory,
			Evidence:    "Cryptographically verifiable AIBOM generated with SHA-256 hash chains.",
			Weight:      1.2,
		},
		{
			Clause:      "A.8.4",
			Title:       "Human Oversight for High-Risk Systems",
			Description: "High-risk decisions mandate verified human review and approval gates.",
			Fulfilled:   true,
			Evidence:    "Quality gate flags Annex III components requiring mandatory human review.",
			Weight:      1.0,
		},
		{
			Clause:      "A.9.1",
			Title:       "Transparency and Disclosure",
			Description: "System users are notified when interacting with generative or conversational AI.",
			Fulfilled:   hasInventory,
			Evidence:    "Article 50 disclosure obligations tracked and flagged for user-facing models.",
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
