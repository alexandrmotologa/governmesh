package remediation

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alexandrmotologa/governmesh/internal/config"
	"github.com/alexandrmotologa/governmesh/internal/scanner"
)

// ProposedFix represents an automated compliance code modification.
type ProposedFix struct {
	FilePath     string `json:"file_path"`
	Line         int    `json:"line"`
	Type         string `json:"type"` // "article_50_disclosure", "human_oversight_gate", "pii_redaction"
	Description  string `json:"description"`
	OriginalText string `json:"original_text"`
	Remediated   string `json:"remediated_text"`
}

// PlanFixes analyzes scan results and generates targeted compliance remediations.
func PlanFixes(res *scanner.ScanResult, baseDir string) ([]ProposedFix, error) {
	var fixes []ProposedFix

	for _, a := range res.Assets {
		fullPath := filepath.Join(baseDir, filepath.FromSlash(a.SourceFilePath))
		contentBytes, err := os.ReadFile(fullPath)
		if err != nil {
			continue
		}
		lines := strings.Split(string(contentBytes), "\n")

		// 1. Article 50 Transparency Disclosure for Chatbots without notice
		if a.Risk.Tier == config.RiskTransparency {
			hasNotice := strings.Contains(strings.ToLower(a.SystemPrompt), "eu ai act") ||
				strings.Contains(strings.ToLower(a.SystemPrompt), "artificial intelligence") ||
				strings.Contains(strings.ToLower(a.SystemPrompt), "ai assistant")

			if !hasNotice && a.LineNumber > 0 && a.LineNumber <= len(lines) {
				origLine := lines[a.LineNumber-1]
				disclosureComment := "  # [EU AI Act Article 50]: Mandatory user disclosure header enforced."
				fixes = append(fixes, ProposedFix{
					FilePath:     a.SourceFilePath,
					Line:         a.LineNumber,
					Type:         "article_50_disclosure",
					Description:  "Inject statutory Article 50 AI disclosure notice into system prompt",
					OriginalText: origLine,
					Remediated:   disclosureComment + "\n" + origLine,
				})
			}
		}

		// 2. Article 14 Human Oversight Gate for High-Risk components
		if a.Risk.Tier == config.RiskHigh && a.Risk.RequiresHumanGate {
			if a.LineNumber > 0 && a.LineNumber <= len(lines) {
				origLine := lines[a.LineNumber-1]
				oversightNote := "  # [EU AI Act Article 14]: Mandatory Human-in-the-Loop review checkpoint required prior to decision enforcement."
				fixes = append(fixes, ProposedFix{
					FilePath:     a.SourceFilePath,
					Line:         a.LineNumber,
					Type:         "human_oversight_gate",
					Description:  "Attach statutory Article 14 human oversight requirement to High-Risk workflow",
					OriginalText: origLine,
					Remediated:   oversightNote + "\n" + origLine,
				})
			}
		}

		// 3. PII masking remediation
		for _, f := range a.PIIFindings {
			if f.Variable != "" && f.Line > 0 && f.Line <= len(lines) {
				origLine := lines[f.Line-1]
				if strings.Contains(origLine, f.Variable) {
					// Propose masking wrapper
					comment := fmt.Sprintf("  # [Privacy Safeguard]: Variable '{%s}' should be tokenized/pseudonymized prior to cloud egress.", f.Variable)
					fixes = append(fixes, ProposedFix{
						FilePath:     a.SourceFilePath,
						Line:         f.Line,
						Type:         "pii_redaction",
						Description:  fmt.Sprintf("Add privacy pseudonymization guard for sensitive variable '%s'", f.Variable),
						OriginalText: origLine,
						Remediated:   comment + "\n" + origLine,
					})
				}
			}
		}
	}

	return fixes, nil
}

// GenerateUnifiedDiff produces standard unified diff format string for git patch consumption.
func GenerateUnifiedDiff(fixes []ProposedFix) string {
	var sb strings.Builder
	fileGroups := make(map[string][]ProposedFix)
	for _, f := range fixes {
		fileGroups[f.FilePath] = append(fileGroups[f.FilePath], f)
	}

	for path, group := range fileGroups {
		sb.WriteString(fmt.Sprintf("--- a/%s\n", path))
		sb.WriteString(fmt.Sprintf("+++ b/%s\n", path))
		for _, fix := range group {
			sb.WriteString(fmt.Sprintf("@@ -%d,1 +%d,2 @@\n", fix.Line, fix.Line))
			sb.WriteString(fmt.Sprintf("-%s\n", fix.OriginalText))
			sb.WriteString(fmt.Sprintf("+%s\n", fix.Remediated))
		}
	}

	return sb.String()
}
