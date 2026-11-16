package report

import (
	"strings"
	"testing"
	"time"

	"github.com/alexandrmotologa/governmesh/internal/attestation"
	"github.com/alexandrmotologa/governmesh/internal/config"
	"github.com/alexandrmotologa/governmesh/internal/scanner"
)

func TestAnnexIVDossierGeneration(t *testing.T) {
	res := &scanner.ScanResult{
		ProjectName:   "FinTech AI Underwriting",
		ScanTimestamp: time.Now(),
		ScannedPath:   "/work/fintech",
		HashChain:     "abcdef1234567890",
		Assets: []scanner.DiscoveredAsset{
			{
				ID:             "fin-credit-1",
				ModelFamily:    "gpt-4o",
				Provider:       "OpenAI",
				Framework:      "openai-sdk",
				SourceFilePath: "credit.py",
				LineNumber:     42,
				Risk: scanner.RiskAssessment{
					Tier:              config.RiskHigh,
					ApplicableArticle: "EU AI Act Annex III Section 5",
					Justification:     "Credit scoring and eligibility evaluation",
					RequiresHumanGate: true,
				},
			},
		},
		RiskSummary:       map[config.RiskLevel]int{config.RiskHigh: 1},
		ComplianceScore:   85.0,
		PassedQualityGate: true,
	}

	att := &attestation.Manifest{
		RootHash:       "roothash123",
		AIBOMHash:      "aibomhash123",
		Signer:         "chief-compliance-officer",
		Organization:   "FinTech Bank",
		Timestamp:      time.Now().UTC().Format(time.RFC3339),
		SignatureBlock: "sigblock123",
	}

	dossier := GenerateAnnexIVDossier(res, att)
	if !strings.Contains(dossier, "EU AI Act Annex IV Technical Documentation") {
		t.Errorf("expected Annex IV header in dossier")
	}
	if !strings.Contains(dossier, "FinTech AI Underwriting") {
		t.Errorf("expected project name in dossier")
	}
	if !strings.Contains(dossier, "fin-credit-1") {
		t.Errorf("expected asset ID in dossier")
	}
	if !strings.Contains(dossier, "roothash123") {
		t.Errorf("expected attestation hash in dossier")
	}
}
