package attestation

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// VerificationResult details the cryptographic validation outcome.
type VerificationResult struct {
	Valid          bool     `json:"valid"`
	ProjectName    string   `json:"project_name"`
	Timestamp      string   `json:"timestamp"`
	Signer         string   `json:"signer"`
	Organization   string   `json:"organization"`
	AIBOMMatches   bool     `json:"aibom_matches"`
	FileTreeValid  bool     `json:"file_tree_valid"`
	SignatureValid bool     `json:"signature_valid"`
	TamperedFiles  []string `json:"tampered_files,omitempty"`
	MissingFiles   []string `json:"missing_files,omitempty"`
	Errors         []string `json:"errors,omitempty"`
}

// VerifyBundle checks the attestation manifest against an AIBOM file and source directory.
func VerifyBundle(manifestPath string, aibomPath string, baseDir string) (*VerificationResult, error) {
	manifestBytes, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("reading attestation manifest: %w", err)
	}

	var m Manifest
	if err := json.Unmarshal(manifestBytes, &m); err != nil {
		return nil, fmt.Errorf("parsing attestation manifest: %w", err)
	}

	result := &VerificationResult{
		Valid:        true,
		ProjectName:  m.ProjectName,
		Timestamp:    m.Timestamp,
		Signer:       m.Signer,
		Organization: m.Organization,
	}

	// 1. Verify AIBOM file hash
	if aibomPath != "" {
		aibomBytes, err := os.ReadFile(aibomPath)
		if err != nil {
			result.AIBOMMatches = false
			result.Valid = false
			result.Errors = append(result.Errors, fmt.Sprintf("failed to read AIBOM file: %v", err))
		} else {
			sum := sha256.Sum256(aibomBytes)
			calcAIBOMHash := hex.EncodeToString(sum[:])
			if calcAIBOMHash == m.AIBOMHash {
				result.AIBOMMatches = true
			} else {
				result.AIBOMMatches = false
				result.Valid = false
				result.Errors = append(result.Errors, fmt.Sprintf("AIBOM hash mismatch: expected %s, got %s", m.AIBOMHash, calcAIBOMHash))
			}
		}
	} else {
		result.AIBOMMatches = true
	}

	// 2. Verify files on disk if baseDir provided
	if baseDir != "" {
		fileTreeValid := true
		for relPath, expectedHash := range m.FileTreeHashes {
			fullPath := filepath.Join(baseDir, filepath.FromSlash(relPath))
			data, err := os.ReadFile(fullPath)
			if err != nil {
				fileTreeValid = false
				result.MissingFiles = append(result.MissingFiles, relPath)
				continue
			}

			sum := sha256.Sum256(data)
			actualHash := hex.EncodeToString(sum[:])
			if actualHash != expectedHash {
				fileTreeValid = false
				result.TamperedFiles = append(result.TamperedFiles, relPath)
			}
		}
		result.FileTreeValid = fileTreeValid
		if !fileTreeValid {
			result.Valid = false
		}
	} else {
		result.FileTreeValid = true
	}

	// 3. Verify signature block consistency
	var sortedFiles []string
	for f := range m.FileTreeHashes {
		sortedFiles = append(sortedFiles, f)
	}
	sort.Strings(sortedFiles)

	hasher := sha256.New()
	for _, f := range sortedFiles {
		hasher.Write([]byte(f + ":" + m.FileTreeHashes[f]))
	}
	expectedRoot := hex.EncodeToString(hasher.Sum(nil))

	if expectedRoot != m.RootHash {
		result.SignatureValid = false
		result.Valid = false
		result.Errors = append(result.Errors, "internal root hash integrity check failed")
	} else {
		sigHasher := sha256.New()
		sigHasher.Write([]byte(m.RootHash + m.AIBOMHash + m.RiskSignature + m.Timestamp + m.Signer))
		calcSig := hex.EncodeToString(sigHasher.Sum(nil))
		if calcSig == m.SignatureBlock {
			result.SignatureValid = true
		} else {
			result.SignatureValid = false
			result.Valid = false
			result.Errors = append(result.Errors, "cryptographic signature block verification failed")
		}
	}

	return result, nil
}
