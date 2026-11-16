package report

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/alexandrmotologa/governmesh/internal/config"
	"github.com/alexandrmotologa/governmesh/internal/pii"
	"github.com/alexandrmotologa/governmesh/internal/scanner"
)

func TestSARIFGeneration(t *testing.T) {
	res := &scanner.ScanResult{
		ProjectName:   "Enterprise SARIF Test",
		ScanTimestamp: time.Now(),
		Assets: []scanner.DiscoveredAsset{
			{
				ID:             "asset-high",
				ModelFamily:    "gpt-4o",
				SourceFilePath: "services/hiring.py",
				LineNumber:     15,
				Risk: scanner.RiskAssessment{
					Tier:              config.RiskHigh,
					ApplicableArticle: "Annex III Section 4",
					Justification:     "Employment resume parsing",
				},
				PIIFindings: []pii.Finding{
					{
						Type:        "ssn",
						Line:        16,
						Severity:    "critical",
						Variable:    "candidate_ssn",
						Description: "National ID variable bound into prompt",
					},
				},
			},
		},
		RiskSummary:       map[config.RiskLevel]int{config.RiskHigh: 1},
		ComplianceScore:   60.0,
		PassedQualityGate: false,
	}

	sarifDoc, err := GenerateSARIF(res)
	if err != nil {
		t.Fatalf("GenerateSARIF failed: %v", err)
	}

	data, err := sarifDoc.ToJSON()
	if err != nil {
		t.Fatalf("SARIF serialization failed: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("invalid SARIF json: %v", err)
	}

	if parsed["version"] != "2.1.0" {
		t.Errorf("expected version 2.1.0, got %v", parsed["version"])
	}

	runs := parsed["runs"].([]interface{})
	if len(runs) == 0 {
		t.Fatalf("expected runs in SARIF")
	}

	firstRun := runs[0].(map[string]interface{})
	results := firstRun["results"].([]interface{})
	if len(results) < 2 {
		t.Errorf("expected at least 2 results (High-Risk + PII finding), got %d", len(results))
	}
}
