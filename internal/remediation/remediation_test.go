package remediation

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alexandrmotologa/governmesh/internal/config"
	"github.com/alexandrmotologa/governmesh/internal/pii"
	"github.com/alexandrmotologa/governmesh/internal/scanner"
)

func TestRemediationPlanning(t *testing.T) {
	tmpDir := t.TempDir()
	sourceFile := filepath.Join(tmpDir, "bot.py")
	content := "def chat_handler(user_input, user_ssn):\n    prompt = f'User SSN: {user_ssn}'\n    return client.chat(prompt)\n"
	if err := os.WriteFile(sourceFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	res := &scanner.ScanResult{
		ProjectName: "Test Project",
		Assets: []scanner.DiscoveredAsset{
			{
				ID:             "bot-1",
				ModelFamily:    "gpt-4o",
				SourceFilePath: "bot.py",
				LineNumber:     1,
				Risk: scanner.RiskAssessment{
					Tier:              config.RiskHigh,
					RequiresHumanGate: true,
				},
				PIIFindings: []pii.Finding{
					{
						Type:     "ssn",
						Line:     2,
						Variable: "user_ssn",
					},
				},
			},
		},
	}

	fixes, err := PlanFixes(res, tmpDir)
	if err != nil {
		t.Fatalf("PlanFixes failed: %v", err)
	}

	if len(fixes) == 0 {
		t.Fatalf("expected proposed fixes, got 0")
	}

	diff := GenerateUnifiedDiff(fixes)
	if !strings.Contains(diff, "--- a/bot.py") {
		t.Errorf("expected unified diff format, got:\n%s", diff)
	}
	if !strings.Contains(diff, "Article 14") {
		t.Errorf("expected Article 14 human oversight in diff, got:\n%s", diff)
	}
	if !strings.Contains(diff, "user_ssn") {
		t.Errorf("expected user_ssn privacy safeguard in diff, got:\n%s", diff)
	}
}
