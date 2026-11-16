package aibom

import (
	"time"

	"github.com/alexandrmotologa/governmesh/internal/config"
	"github.com/alexandrmotologa/governmesh/internal/scanner"
)

// AIBOMMetadata contains generation parameters and cryptographic attestation links.
type AIBOMMetadata struct {
	SpecVersion     string           `json:"spec_version"`
	Standard        string           `json:"standard"` // "CycloneDX-1.6" or "SPDX-3.0"
	GeneratedAt     time.Time        `json:"generated_at"`
	ProjectName     string           `json:"project_name"`
	ToolVersion     string           `json:"tool_version"`
	HashChain       string           `json:"hash_chain"`
	ComplianceScore float64          `json:"compliance_score"`
	OverallRisk     config.RiskLevel `json:"overall_risk"`
}

// ModelCardSummary holds summarized governance characteristics.
type ModelCardSummary struct {
	ModelFamily     string   `json:"model_family"`
	Provider        string   `json:"provider"`
	Framework       string   `json:"framework"`
	RiskTier        string   `json:"risk_tier"`
	ApplicableRule  string   `json:"applicable_rule"`
	EgressTarget    string   `json:"egress_target"`
	PIIExposureRisk string   `json:"pii_exposure_risk"`
	SourceLocation  string   `json:"source_location"`
	NISTCategories  []string `json:"nist_categories"`
	ISO42001        []string `json:"iso_42001_controls"`
}

// ConvertScanToSummaries converts scanner assets into clean AIBOM model card records.
func ConvertScanToSummaries(assets []scanner.DiscoveredAsset) []ModelCardSummary {
	var list []ModelCardSummary
	for _, a := range assets {
		piiRisk := "none"
		if len(a.PIIFindings) > 0 {
			piiRisk = "detected"
			for _, f := range a.PIIFindings {
				if f.Severity == "critical" {
					piiRisk = "critical"
					break
				}
			}
		}

		list = append(list, ModelCardSummary{
			ModelFamily:     a.ModelFamily,
			Provider:        a.Provider,
			Framework:       a.Framework,
			RiskTier:        string(a.Risk.Tier),
			ApplicableRule:  a.Risk.ApplicableArticle,
			EgressTarget:    a.EgressTarget,
			PIIExposureRisk: piiRisk,
			SourceLocation:  a.SourceFilePath,
			NISTCategories:  a.Risk.NISTCategories,
			ISO42001:        a.Risk.ISO42001Controls,
		})
	}
	return list
}
