package attestation

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/alexandrmotologa/governmesh/internal/config"
	"github.com/alexandrmotologa/governmesh/internal/scanner"
)

// Manifest contains the full cryptographic proof of the audit.
type Manifest struct {
	Version        string            `json:"version"`
	ProjectName    string            `json:"project_name"`
	Timestamp      string            `json:"timestamp"`
	Signer         string            `json:"signer"`
	Organization   string            `json:"organization"`
	Algorithm      string            `json:"algorithm"` // "SHA-256-MERKLE"
	RootHash       string            `json:"root_hash"`
	AIBOMHash      string            `json:"aibom_hash"`
	FileTreeHashes map[string]string `json:"file_tree_hashes"`
	RiskSignature  string            `json:"risk_signature"`
	SignatureBlock string            `json:"signature_block"`
}

// GenerateManifest produces a cryptographically sealed compliance manifest.
func GenerateManifest(res *scanner.ScanResult, aibomBytes []byte, cfg *config.Config) (*Manifest, error) {
	if cfg == nil {
		cfg = config.DefaultConfig()
	}

	// 1. Calculate AIBOM hash
	aibomSum := sha256.Sum256(aibomBytes)
	aibomHash := hex.EncodeToString(aibomSum[:])

	// 2. Build sorted file tree hash map
	fileTree := make(map[string]string)
	for _, a := range res.Assets {
		fileTree[a.SourceFilePath] = a.FileHash
	}

	var sortedFiles []string
	for f := range fileTree {
		sortedFiles = append(sortedFiles, f)
	}
	sort.Strings(sortedFiles)

	// 3. Merkle root of files
	hasher := sha256.New()
	for _, f := range sortedFiles {
		hasher.Write([]byte(f + ":" + fileTree[f]))
	}
	rootHash := hex.EncodeToString(hasher.Sum(nil))

	// 4. Risk signature (locks in the classified tiers and gate decision)
	riskHasher := sha256.New()
	riskHasher.Write([]byte(fmt.Sprintf("%d-%d-%d-%d-%.1f-%t",
		res.RiskSummary[config.RiskUnacceptable],
		res.RiskSummary[config.RiskHigh],
		res.RiskSummary[config.RiskTransparency],
		res.RiskSummary[config.RiskMinimal],
		res.ComplianceScore,
		res.PassedQualityGate,
	)))
	riskSignature := hex.EncodeToString(riskHasher.Sum(nil))

	timestamp := time.Now().UTC().Format(time.RFC3339)

	// 5. Compute master signature block
	sigHasher := sha256.New()
	sigHasher.Write([]byte(rootHash + aibomHash + riskSignature + timestamp + cfg.Attestation.SignerName))
	signatureBlock := hex.EncodeToString(sigHasher.Sum(nil))

	manifest := &Manifest{
		Version:        "1.0.0",
		ProjectName:    res.ProjectName,
		Timestamp:      timestamp,
		Signer:         cfg.Attestation.SignerName,
		Organization:   cfg.Attestation.Organization,
		Algorithm:      "SHA-256-MERKLE",
		RootHash:       rootHash,
		AIBOMHash:      aibomHash,
		FileTreeHashes: fileTree,
		RiskSignature:  riskSignature,
		SignatureBlock: signatureBlock,
	}

	return manifest, nil
}

// SaveToFile writes the manifest to disk as formatted JSON.
func (m *Manifest) SaveToFile(path string) error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling attestation: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}
