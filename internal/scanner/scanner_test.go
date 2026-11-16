package scanner

import (
	"path/filepath"
	"testing"

	"github.com/alexandrmotologa/governmesh/internal/config"
)

func TestScannerEndToEnd(t *testing.T) {
	sampleDir, err := filepath.Abs("../../examples/enterprise_ai_service")
	if err != nil {
		t.Fatal(err)
	}

	cfg := config.DefaultConfig()
	scn := NewScanner(cfg, nil)

	res, err := scn.Run(sampleDir)
	if err != nil {
		t.Fatalf("scanner error: %v", err)
	}

	if len(res.Assets) == 0 {
		t.Fatalf("expected discovered AI assets, got 0")
	}

	foundPython := false
	foundTS := false
	foundManifest := false
	foundPrompt := false

	for _, a := range res.Assets {
		switch a.Language {
		case "python":
			foundPython = true
		case "typescript":
			foundTS = true
		case "manifest":
			foundManifest = true
		case "prompt":
			foundPrompt = true
		}
	}

	if !foundPython {
		t.Errorf("expected to find Python AI integrations")
	}
	if !foundTS {
		t.Errorf("expected to find TypeScript AI integrations")
	}
	if !foundManifest {
		t.Errorf("expected to find Docker/Manifest AI integrations")
	}
	if !foundPrompt {
		t.Errorf("expected to find Prompt template integrations")
	}

	if res.HashChain == "" {
		t.Errorf("expected non-empty hash chain")
	}
}
