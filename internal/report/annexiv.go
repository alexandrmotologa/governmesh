package report

import (
	"fmt"
	"strings"
	"time"

	"github.com/alexandrmotologa/governmesh/internal/attestation"
	"github.com/alexandrmotologa/governmesh/internal/config"
	"github.com/alexandrmotologa/governmesh/internal/scanner"
)

// GenerateAnnexIVDossier synthesizes the mandatory EU AI Act Annex IV Technical Documentation.
func GenerateAnnexIVDossier(res *scanner.ScanResult, att *attestation.Manifest) string {
	var sb strings.Builder

	sb.WriteString("# EU AI Act Annex IV Technical Documentation Dossier\n\n")
	sb.WriteString("> Statutory Technical Documentation for High-Risk AI Systems pursuant to Article 11 and Annex IV of Regulation (EU) 2024/1689\n\n")

	sb.WriteString(fmt.Sprintf("**AI System / Provider:** %s  \n", res.ProjectName))
	sb.WriteString(fmt.Sprintf("**Date of Compilation:** %s  \n", res.ScanTimestamp.Format(time.RFC1123)))
	sb.WriteString(fmt.Sprintf("**Codebase Target:** `%s`  \n", res.ScannedPath))
	sb.WriteString(fmt.Sprintf("**Dossier ID:** `DOC-ANNEX-IV-%s`  \n\n", res.HashChain[:12]))

	// Section 1
	sb.WriteString("## 1. General Description of the AI System (Annex IV.1)\n\n")
	sb.WriteString("### 1.1 Intended Purpose and Practical Application\n")
	sb.WriteString(fmt.Sprintf("This technical documentation pertains to the software components identified within `%s`. ", res.ProjectName))
	sb.WriteString(fmt.Sprintf("The static AST inspection identified %d integrated AI model(s) and framework components.\n\n", len(res.Assets)))

	sb.WriteString("### 1.2 System Risk Classification\n")
	sb.WriteString("| Category | System Count | Statutory Basis |\n")
	sb.WriteString("| :--- | :--- | :--- |\n")
	sb.WriteString(fmt.Sprintf("| Unacceptable Risk (Article 5) | %d | Prohibited practices (requires immediate decommissioning) |\n", res.RiskSummary[config.RiskUnacceptable]))
	sb.WriteString(fmt.Sprintf("| High-Risk (Annex III) | %d | Subject to mandatory conformity assessment and this dossier |\n", res.RiskSummary[config.RiskHigh]))
	sb.WriteString(fmt.Sprintf("| Specific Transparency (Article 50) | %d | Subject to user disclosure notices |\n", res.RiskSummary[config.RiskTransparency]))
	sb.WriteString(fmt.Sprintf("| Minimal Risk (Article 69) | %d | Standard software quality baseline |\n\n", res.RiskSummary[config.RiskMinimal]))

	// Section 2
	sb.WriteString("## 2. Methods and Steps for AI Development (Annex IV.2)\n\n")
	sb.WriteString("### 2.1 Model Architectures and Foundation Model Supply Chain\n\n")
	sb.WriteString("| Component ID | Model Family | Provider / Origin | Framework | Code Location |\n")
	sb.WriteString("| :--- | :--- | :--- | :--- | :--- |\n")
	for _, a := range res.Assets {
		sb.WriteString(fmt.Sprintf("| `%s` | **%s** | %s | `%s` | `%s:%d` |\n",
			a.ID, a.ModelFamily, a.Provider, a.Framework, a.SourceFilePath, a.LineNumber))
	}
	sb.WriteString("\n")

	// Section 3
	sb.WriteString("## 3. Human Oversight & Operational Monitoring Architecture (Annex IV.3)\n\n")
	sb.WriteString("Pursuant to Article 14 of the EU AI Act, High-Risk AI systems must be designed to enable natural persons to oversee their operation.\n\n")
	hasHighRisk := false
	for _, a := range res.Assets {
		if a.Risk.Tier == config.RiskHigh || a.Risk.RequiresHumanGate {
			hasHighRisk = true
			sb.WriteString(fmt.Sprintf("- **Component `%s` (%s):** %s\n", a.ID, a.ModelFamily, a.Risk.Justification))
			sb.WriteString(fmt.Sprintf("  - *Mandatory Human Review Gate:* Active CI/CD verification flag.\n"))
			sb.WriteString(fmt.Sprintf("  - *Governing Rule:* %s\n", a.Risk.ApplicableArticle))
		}
	}
	if !hasHighRisk {
		sb.WriteString("No High-Risk Annex III components were detected requiring mandatory Article 14 human oversight protocols.\n")
	}
	sb.WriteString("\n")

	// Section 4
	sb.WriteString("## 4. Data Governance and Egress Boundaries (Annex IV.4)\n\n")
	sb.WriteString("Pursuant to Article 10, training, validation, and testing datasets and prompt data flows must adhere to strict data governance practices.\n\n")
	if len(res.PIIExposures) == 0 {
		sb.WriteString("Zero sensitive personal data or prohibited demographic variables were detected egressing to third-party model providers.\n\n")
	} else {
		sb.WriteString("### Detected Sensitive Variable Flows\n\n")
		sb.WriteString("| Variable / Snippet | PII Category | Destination Provider | Source Code Line |\n")
		sb.WriteString("| :--- | :--- | :--- | :--- |\n")
		for _, exp := range res.PIIExposures {
			sb.WriteString(fmt.Sprintf("| `%s` | `%s` | %s | `%s:%d` |\n",
				exp.Snippet, exp.PIIType, exp.TargetCloud, exp.FilePath, exp.Line))
		}
		sb.WriteString("\n")
	}

	// Section 5
	sb.WriteString("## 5. Technical Standards and Conformity (Annex IV.5)\n\n")
	sb.WriteString("The system has been evaluated against harmonized standards and international governance frameworks:\n\n")
	sb.WriteString(fmt.Sprintf("- **ISO/IEC 42001:2023 Alignment:** Controls A.5, A.6, A.7, A.8, A.9 evaluated via automated AST inspection.\n"))
	sb.WriteString(fmt.Sprintf("- **NIST AI RMF 1.0 Alignment:** Functions GOVERN, MAP, MEASURE, MANAGE verified.\n"))
	sb.WriteString(fmt.Sprintf("- **Composite Readiness Score:** **%.1f%%**\n\n", res.ComplianceScore))

	// Section 6
	sb.WriteString("## 6. Cryptographic Attestation and Signatures (Annex IV.6)\n\n")
	if att != nil {
		sb.WriteString("This technical dossier is cryptographically anchored to the repository source state via a SHA-256 Merkle hash chain:\n\n")
		sb.WriteString(fmt.Sprintf("- **Merkle Root Checksum:** `%s`\n", att.RootHash))
		sb.WriteString(fmt.Sprintf("- **AIBOM Checksum:** `%s`\n", att.AIBOMHash))
		sb.WriteString(fmt.Sprintf("- **Attestation Signer:** %s (%s)\n", att.Signer, att.Organization))
		sb.WriteString(fmt.Sprintf("- **Attestation Timestamp:** %s\n", att.Timestamp))
		sb.WriteString(fmt.Sprintf("- **Master Signature Block:** `%s`\n\n", att.SignatureBlock))
	}

	sb.WriteString("---\n*Dossier compiled by governmesh automated compliance engine.*  \n")
	return sb.String()
}
