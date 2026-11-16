package aibom

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/alexandrmotologa/governmesh/internal/config"
	"github.com/alexandrmotologa/governmesh/internal/scanner"
)

func TestCycloneDXAndSPDXGeneration(t *testing.T) {
	res := &scanner.ScanResult{
		ProjectName:   "Enterprise Test",
		ScanTimestamp: time.Now(),
		Assets: []scanner.DiscoveredAsset{
			{
				ID:             "asset-1",
				ModelFamily:    "gpt-4o",
				Provider:       "OpenAI",
				Framework:      "openai-sdk",
				SourceFilePath: "app.py",
				LineNumber:     10,
				EgressTarget:   "external_cloud",
				Risk: scanner.RiskAssessment{
					Tier:              config.RiskHigh,
					ApplicableArticle: "Annex III Section 4",
					Justification:     "Employment CV screening",
				},
				FileHash: "abc123hash",
			},
		},
		RiskSummary:       map[config.RiskLevel]int{config.RiskHigh: 1},
		ComplianceScore:   85.0,
		PassedQualityGate: true,
		HashChain:         "chain123",
	}

	// 1. CycloneDX 1.6
	cdx, err := GenerateCycloneDX(res)
	if err != nil {
		t.Fatalf("CycloneDX generation failed: %v", err)
	}
	cdxBytes, err := cdx.ToJSON()
	if err != nil {
		t.Fatalf("CycloneDX serialization failed: %v", err)
	}

	var parsedCDX map[string]interface{}
	if err := json.Unmarshal(cdxBytes, &parsedCDX); err != nil {
		t.Fatalf("CycloneDX invalid json: %v", err)
	}
	if parsedCDX["bomFormat"] != "CycloneDX" || parsedCDX["specVersion"] != "1.6" {
		t.Errorf("expected CycloneDX 1.6 header")
	}

	// 2. SPDX 3.0
	spdxDoc, err := GenerateSPDX(res)
	if err != nil {
		t.Fatalf("SPDX generation failed: %v", err)
	}
	spdxBytes, err := spdxDoc.ToJSON()
	if err != nil {
		t.Fatalf("SPDX serialization failed: %v", err)
	}

	var parsedSPDX map[string]interface{}
	if err := json.Unmarshal(spdxBytes, &parsedSPDX); err != nil {
		t.Fatalf("SPDX invalid json: %v", err)
	}
	if parsedSPDX["spdxVersion"] != "SPDX-3.0" {
		t.Errorf("expected SPDX-3.0 header")
	}
}
