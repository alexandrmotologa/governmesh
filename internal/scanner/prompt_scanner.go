package scanner

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/alexandrmotologa/governmesh/internal/pii"
)

var (
	promptVarRe = regexp.MustCompile(`\{([a-zA-Z0-9_\.]+)\}`)
)

// ScanPromptFile inspects standalone prompt files or YAML template definitions.
func ScanPromptFile(path string, baseDir string) ([]DiscoveredAsset, error) {
	contentBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading prompt file: %w", err)
	}

	content := string(contentBytes)
	lowerName := strings.ToLower(filepath.Base(path))

	isPromptFile := strings.HasSuffix(lowerName, ".prompt") ||
		strings.Contains(lowerName, "prompt") ||
		strings.Contains(lowerName, "system_message") ||
		strings.Contains(lowerName, "few_shot")

	if !isPromptFile {
		return nil, nil
	}

	relPath, err := filepath.Rel(baseDir, path)
	if err != nil {
		relPath = path
	}

	hashSum := sha256.Sum256(contentBytes)
	fileHash := hex.EncodeToString(hashSum[:])

	// Extract template variables
	var promptVars []string
	varSeen := make(map[string]bool)
	matches := promptVarRe.FindAllStringSubmatch(content, -1)
	for _, m := range matches {
		if len(m) > 1 && !varSeen[m[1]] {
			varSeen[m[1]] = true
			promptVars = append(promptVars, m[1])
		}
	}

	// Scan literal PII and variable bindings
	piiFindings := pii.ScanPromptText(content, 1)
	piiFindings = append(piiFindings, pii.ScanVariables(promptVars, 1)...)

	lines := strings.Split(content, "\n")
	previewSnippet := content
	if len(lines) > 6 {
		previewSnippet = strings.Join(lines[:6], "\n") + "\n..."
	}

	asset := DiscoveredAsset{
		ID:             fmt.Sprintf("prompt-%s", filepath.Base(relPath)),
		Name:           fmt.Sprintf("Prompt Template: %s", filepath.Base(relPath)),
		ModelFamily:    "Prompt Specification",
		Provider:       "Prompt Engineering Asset",
		SourceFilePath: filepath.ToSlash(relPath),
		LineNumber:     1,
		Language:       "prompt",
		Framework:      "prompt-template",
		ContextSnippet: previewSnippet,
		SystemPrompt:   content,
		PromptVars:     promptVars,
		EgressTarget:   "external_cloud",
		PIIFindings:    piiFindings,
		FileHash:       fileHash,
	}

	return []DiscoveredAsset{asset}, nil
}
