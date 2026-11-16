package report

import (
	"fmt"
	"strings"
	"time"

	"github.com/alexandrmotologa/governmesh/internal/attestation"
	"github.com/alexandrmotologa/governmesh/internal/config"
	"github.com/alexandrmotologa/governmesh/internal/scanner"
)

// GenerateMarkdownReport constructs an executive compliance audit document in Markdown.
func GenerateMarkdownReport(res *scanner.ScanResult, att *attestation.Manifest) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# AI Governance & Regulatory Compliance Audit Report\n\n"))
	sb.WriteString(fmt.Sprintf("**Target Entity / Project:** %s  \n", res.ProjectName))
	sb.WriteString(fmt.Sprintf("**Audit Date:** %s  \n", res.ScanTimestamp.Format(time.RFC1123)))
	sb.WriteString(fmt.Sprintf("**Scanned Codebase Path:** `%s`  \n", res.ScannedPath))
	sb.WriteString(fmt.Sprintf("**Audit Engine:** governmesh v1.0.0 (Go Engine)  \n\n"))

	// Status badge
	gateStatus := "PASSED"
	if !res.PassedQualityGate {
		gateStatus = "FAILED (Non-Compliant Risks Detected)"
	}
	sb.WriteString(fmt.Sprintf("### Executive Status: **%s**\n", gateStatus))
	sb.WriteString(fmt.Sprintf("- **Composite Compliance Score:** **%.1f%%**\n", res.ComplianceScore))
	sb.WriteString(fmt.Sprintf("- **Total Files Inspected:** %d\n", res.TotalFiles))
	sb.WriteString(fmt.Sprintf("- **AI Systems & Integrations Identified:** %d\n\n", len(res.Assets)))

	if len(res.GateFailures) > 0 {
		sb.WriteString("#### Critical Compliance Gate Violations\n")
		for _, f := range res.GateFailures {
			sb.WriteString(fmt.Sprintf("- %s\n", f))
		}
		sb.WriteString("\n")
	}

	// 1. EU AI Act Risk Distribution
	sb.WriteString("## 1. EU AI Act Risk Classification\n\n")
	sb.WriteString("| Risk Tier | Statutory Basis | Systems Count | CI/CD Action Mandate |\n")
	sb.WriteString("| :--- | :--- | :--- | :--- |\n")
	sb.WriteString(fmt.Sprintf("| **Unacceptable Risk** | Article 5 (Prohibited Practices) | %d | Hard CI/CD Gate Block |\n", res.RiskSummary[config.RiskUnacceptable]))
	sb.WriteString(fmt.Sprintf("| **High-Risk** | Article 6 & Annex III | %d | Mandatory Conformity & Human Oversight |\n", res.RiskSummary[config.RiskHigh]))
	sb.WriteString(fmt.Sprintf("| **Specific Transparency** | Article 50 (Chatbots/Synthetic) | %d | End-User AI Disclosure Badge |\n", res.RiskSummary[config.RiskTransparency]))
	sb.WriteString(fmt.Sprintf("| **Minimal Risk** | Article 69 (General Purpose) | %d | Permitted with AIBOM tracking |\n\n", res.RiskSummary[config.RiskMinimal]))

	// 2. Inventory of Discovered AI Assets
	sb.WriteString("## 2. AI Bill of Materials (AIBOM) Inventory\n\n")
	sb.WriteString("| Model / Component | Provider | Framework | Risk Tier | Egress Target | Source File |\n")
	sb.WriteString("| :--- | :--- | :--- | :--- | :--- | :--- |\n")
	for _, a := range res.Assets {
		sb.WriteString(fmt.Sprintf("| `%s` | %s | %s | **%s** | %s | `%s:%d` |\n",
			a.ModelFamily, a.Provider, a.Framework, a.Risk.Tier, a.EgressTarget, a.SourceFilePath, a.LineNumber))
	}
	sb.WriteString("\n")

	// 3. PII & Data Egress Audit
	sb.WriteString("## 3. Sensitive Data Egress & PII Findings\n\n")
	if len(res.PIIExposures) == 0 {
		sb.WriteString("No sensitive personal data (PII/PHI) was detected egressing to external cloud AI models.\n\n")
	} else {
		sb.WriteString("| File | Line | PII Type | Cloud Provider | Redacted Sample | Prohibited |\n")
		sb.WriteString("| :--- | :--- | :--- | :--- | :--- | :--- |\n")
		for _, p := range res.PIIExposures {
			proh := "No"
			if p.IsProhibited {
				proh = "**YES**"
			}
			sb.WriteString(fmt.Sprintf("| `%s` | %d | `%s` | %s | `%s` | %s |\n",
				p.FilePath, p.Line, p.PIIType, p.TargetCloud, p.Snippet, proh))
		}
		sb.WriteString("\n")
	}

	// 4. Cryptographic Attestation Block
	if att != nil {
		sb.WriteString("## 4. Cryptographic Attestation & Audit Trail\n\n")
		sb.WriteString("```json\n{\n")
		sb.WriteString(fmt.Sprintf("  \"algorithm\": \"%s\",\n", att.Algorithm))
		sb.WriteString(fmt.Sprintf("  \"signer\": \"%s\",\n", att.Signer))
		sb.WriteString(fmt.Sprintf("  \"organization\": \"%s\",\n", att.Organization))
		sb.WriteString(fmt.Sprintf("  \"timestamp\": \"%s\",\n", att.Timestamp))
		sb.WriteString(fmt.Sprintf("  \"merkle_root_hash\": \"%s\",\n", att.RootHash))
		sb.WriteString(fmt.Sprintf("  \"aibom_hash\": \"%s\",\n", att.AIBOMHash))
		sb.WriteString(fmt.Sprintf("  \"signature_block\": \"%s\"\n", att.SignatureBlock))
		sb.WriteString("}\n```\n")
	}

	return sb.String()
}
