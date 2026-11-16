package classifier

import (
	"testing"

	"github.com/alexandrmotologa/governmesh/internal/config"
	"github.com/alexandrmotologa/governmesh/internal/scanner"
)

func TestEUAIActClassification(t *testing.T) {
	tests := []struct {
		name         string
		asset        scanner.DiscoveredAsset
		expectedTier config.RiskLevel
	}{
		{
			name: "Prohibited Social Scoring",
			asset: scanner.DiscoveredAsset{
				Name:           "citizen_score_evaluator",
				ModelFamily:    "gpt-4o",
				IntentKeywords: []string{"social_scoring", "citizen_rating"},
			},
			expectedTier: config.RiskUnacceptable,
		},
		{
			name: "High-Risk Recruitment Screening",
			asset: scanner.DiscoveredAsset{
				Name:           "evaluate_candidate_cv",
				ModelFamily:    "gpt-4o",
				IntentKeywords: []string{"evaluate_candidate_cv", "resume_screening"},
			},
			expectedTier: config.RiskHigh,
		},
		{
			name: "High-Risk Credit Underwriting",
			asset: scanner.DiscoveredAsset{
				Name:           "credit_underwriting_engine",
				ModelFamily:    "gpt-4o",
				IntentKeywords: []string{"credit_scoring", "loan_underwriting"},
			},
			expectedTier: config.RiskHigh,
		},
		{
			name: "Article 50 Conversational Assistant",
			asset: scanner.DiscoveredAsset{
				Name:        "support_assistant",
				ModelFamily: "gpt-4o-mini",
			},
			expectedTier: config.RiskTransparency,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assessment := ClassifyAsset(&tc.asset)
			if assessment.Tier != tc.expectedTier {
				t.Errorf("expected tier %s, got %s", tc.expectedTier, assessment.Tier)
			}
		})
	}
}
