package attestation

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/alexandrmotologa/governmesh/internal/config"
	"github.com/alexandrmotologa/governmesh/internal/scanner"
)

func TestAttestationLifecycle(t *testing.T) {
	tmpDir := t.TempDir()
	sampleFile := filepath.Join(tmpDir, "model.py")
	if err := os.WriteFile(sampleFile, []byte("import openai"), 0644); err != nil {
		t.Fatal(err)
	}

	res := &scanner.ScanResult{
		ProjectName:   "Test AI Project",
		ScanTimestamp: time.Now(),
		Assets: []scanner.DiscoveredAsset{
			{
				ID:             "asset-1",
				ModelFamily:    "gpt-4o",
				SourceFilePath: "model.py",
				FileHash:       "dummyhash",
			},
		},
		RiskSummary:       map[config.RiskLevel]int{config.RiskMinimal: 1},
		ComplianceScore:   95.0,
		PassedQualityGate: true,
	}

	manifest, err := GenerateManifest(res, []byte(`{"bomFormat":"CycloneDX"}`), config.DefaultConfig())
	if err != nil {
		t.Fatalf("failed to generate manifest: %v", err)
	}

	manifestPath := filepath.Join(tmpDir, "attestation.json")
	if err := manifest.SaveToFile(manifestPath); err != nil {
		t.Fatalf("failed to save manifest: %v", err)
	}

	aibomPath := filepath.Join(tmpDir, "aibom.json")
	if err := os.WriteFile(aibomPath, []byte(`{"bomFormat":"CycloneDX"}`), 0644); err != nil {
		t.Fatal(err)
	}

	// Verify intact
	result, err := VerifyBundle(manifestPath, aibomPath, "")
	if err != nil {
		t.Fatalf("verification error: %v", err)
	}
	if !result.Valid {
		t.Errorf("expected bundle to be valid, errors: %v", result.Errors)
	}

	// Verify tampering detection
	badAIBOM := filepath.Join(tmpDir, "bad_aibom.json")
	if err := os.WriteFile(badAIBOM, []byte(`{"tampered":true}`), 0644); err != nil {
		t.Fatal(err)
	}
	badResult, err := VerifyBundle(manifestPath, badAIBOM, "")
	if err != nil {
		t.Fatalf("verification error on tampered: %v", err)
	}
	if badResult.Valid {
		t.Errorf("expected tampered bundle to be invalid")
	}
}
