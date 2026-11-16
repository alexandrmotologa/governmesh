package classifier

import (
	"fmt"
	"strings"

	"github.com/alexandrmotologa/governmesh/internal/config"
	"github.com/alexandrmotologa/governmesh/internal/pii"
	"github.com/alexandrmotologa/governmesh/internal/scanner"
)

// ProcessScanResult runs all regulatory classifiers, NIST/ISO scoring, and quality gates on the scan result.
func ProcessScanResult(res *scanner.ScanResult, cfg *config.Config) {
	if cfg == nil {
		cfg = config.DefaultConfig()
	}

	riskCounts := map[config.RiskLevel]int{
		config.RiskUnacceptable: 0,
		config.RiskHigh:         0,
		config.RiskTransparency: 0,
		config.RiskMinimal:      0,
	}

	var piiExposures []pii.Exposure
	var allPIIFindings []pii.Finding

	for i := range res.Assets {
		asset := &res.Assets[i]
		asset.Risk = ClassifyAsset(asset)
		riskCounts[asset.Risk.Tier]++

		for _, f := range asset.PIIFindings {
			allPIIFindings = append(allPIIFindings, f)
			if asset.EgressTarget == "external_cloud" {
				isProhibited := f.Severity == "critical" || asset.Risk.Tier == config.RiskUnacceptable
				piiExposures = append(piiExposures, pii.Exposure{
					AssetID:      asset.ID,
					FilePath:     asset.SourceFilePath,
					Line:         f.Line,
					PIIType:      f.Type,
					TargetCloud:  asset.Provider,
					Snippet:      f.Snippet,
					IsProhibited: isProhibited,
				})
			}
		}
	}

	res.RiskSummary = riskCounts
	res.PIIExposures = piiExposures

	_, nistScore := EvaluateNIST(res.Assets, allPIIFindings)
	_, isoScore := EvaluateISO42001(res.Assets, allPIIFindings)

	composite := (nistScore + isoScore) / 2.0
	if riskCounts[config.RiskUnacceptable] > 0 {
		composite -= float64(riskCounts[config.RiskUnacceptable] * 35)
	}
	if riskCounts[config.RiskHigh] > 0 {
		composite -= float64(riskCounts[config.RiskHigh] * 10)
	}
	if len(piiExposures) > 0 {
		composite -= float64(len(piiExposures) * 5)
	}
	if composite < 0 {
		composite = 0
	}
	res.ComplianceScore = float64(int(composite*10)) / 10.0

	gatePassed := true
	var failures []string

	if riskCounts[config.RiskUnacceptable] > 0 {
		gatePassed = false
		failures = append(failures, fmt.Sprintf(
			"GATE FAILED: %d AI system(s) flagged as EU AI Act Article 5 Prohibited / Unacceptable Risk.",
			riskCounts[config.RiskUnacceptable],
		))
	}

	if cfg.Compliance.FailOnRiskLevel == config.RiskHigh && riskCounts[config.RiskHigh] > 0 {
		gatePassed = false
		failures = append(failures, fmt.Sprintf(
			"GATE FAILED: %d AI system(s) classified as High-Risk (Annex III) exceeding threshold.",
			riskCounts[config.RiskHigh],
		))
	}

	if cfg.Compliance.FailOnPIIEgress && len(piiExposures) > 0 {
		var critCount int
		for _, exp := range piiExposures {
			if exp.IsProhibited {
				critCount++
			}
		}
		if critCount > 0 {
			gatePassed = false
			failures = append(failures, fmt.Sprintf(
				"GATE FAILED: %d critical PII/PHI variable(s) egressing to external cloud LLM providers.",
				critCount,
			))
		}
	}

	res.PassedQualityGate = gatePassed
	res.GateFailures = failures
}

// SummaryText formats a terminal friendly summary.
func SummaryText(res *scanner.ScanResult) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Scan completed in %d ms across %d files (%d AI assets identified).\n", res.DurationMs, res.TotalFiles, len(res.Assets)))
	sb.WriteString(fmt.Sprintf("EU AI Act Classification:\n"))
	sb.WriteString(fmt.Sprintf("  - Unacceptable Risk: %d\n", res.RiskSummary[config.RiskUnacceptable]))
	sb.WriteString(fmt.Sprintf("  - High-Risk (Annex III): %d\n", res.RiskSummary[config.RiskHigh]))
	sb.WriteString(fmt.Sprintf("  - Specific Transparency (Art 50): %d\n", res.RiskSummary[config.RiskTransparency]))
	sb.WriteString(fmt.Sprintf("  - Minimal Risk: %d\n", res.RiskSummary[config.RiskMinimal]))
	sb.WriteString(fmt.Sprintf("Compliance Score: %.1f%%\n", res.ComplianceScore))
	if res.PassedQualityGate {
		sb.WriteString("CI/CD Quality Gate: PASSED\n")
	} else {
		sb.WriteString("CI/CD Quality Gate: FAILED\n")
		for _, f := range res.GateFailures {
			sb.WriteString(fmt.Sprintf("  * %s\n", f))
		}
	}
	return sb.String()
}
