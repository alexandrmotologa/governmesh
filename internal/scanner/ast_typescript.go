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
	tsAIImportRe = regexp.MustCompile(`(?m)^\s*(?:import\s+.*?\s+from\s+["'](@ai-sdk/[a-zA-Z0-9_-]+|ai|openai|@anthropic-ai/sdk|@langchain/[a-zA-Z0-9_-]+|langchain|@google/genai)["']|const\s+.*?\s*=\s*require\(["'](openai|@anthropic-ai/sdk|ai|@langchain/[a-zA-Z0-9_-]+)["']\))`)
	tsModelRefRe = regexp.MustCompile(`(?:model:\s*(?:openai|anthropic|google|cohere)\(["']([a-zA-Z0-9\._\-\/]+)["']\)|model:\s*["']([a-zA-Z0-9\._\-\/]+)["']|modelId:\s*["']([a-zA-Z0-9\._\-\/]+)["'])`)
	tsSystemPromptRe = regexp.MustCompile("(?s)system:\\s*[\"'`]{1,3}(.*?)[\"'`]{1,3}")
	tsTemplateVarRe = regexp.MustCompile(`\$\{([a-zA-Z0-9_\.]+)\}`)
	tsFunctionDefRe = regexp.MustCompile(`(?:function\s+([a-zA-Z0-9_]+)|const\s+([a-zA-Z0-9_]+)\s*=\s*(?:async\s*)?\()`)
)

// ScanTypeScriptFile inspects a TS/JS file for AI SDK usage and prompt interpolations.
func ScanTypeScriptFile(path string, baseDir string) ([]DiscoveredAsset, error) {
	contentBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading ts/js file: %w", err)
	}

	content := string(contentBytes)
	hasAIImport := tsAIImportRe.MatchString(content)
	hasModelRef := tsModelRefRe.MatchString(content)

	if !hasAIImport && !hasModelRef {
		return nil, nil
	}

	relPath, err := filepath.Rel(baseDir, path)
	if err != nil {
		relPath = path
	}

	hashSum := sha256.Sum256(contentBytes)
	fileHash := hex.EncodeToString(hashSum[:])

	lines := strings.Split(content, "\n")
	var assets []DiscoveredAsset

	// Intent keywords from function definitions
	var intentKeywords []string
	fnMatches := tsFunctionDefRe.FindAllStringSubmatch(content, -1)
	for _, m := range fnMatches {
		if len(m) > 1 && m[1] != "" {
			intentKeywords = append(intentKeywords, strings.ToLower(m[1]))
		} else if len(m) > 2 && m[2] != "" {
			intentKeywords = append(intentKeywords, strings.ToLower(m[2]))
		}
	}

	// Model extraction
	modelMatches := tsModelRefRe.FindAllStringSubmatch(content, -1)
	modelsFound := make(map[string]bool)
	for _, m := range modelMatches {
		for i := 1; i < len(m); i++ {
			if m[i] != "" {
				modelsFound[m[i]] = true
			}
		}
	}

	// System prompt
	systemPrompt := ""
	if sm := tsSystemPromptRe.FindStringSubmatch(content); len(sm) > 1 {
		systemPrompt = strings.TrimSpace(sm[1])
	}

	// Template interpolation variables
	var promptVars []string
	varSeen := make(map[string]bool)
	varMatches := tsTemplateVarRe.FindAllStringSubmatch(content, -1)
	for _, v := range varMatches {
		if len(v) > 1 && !varSeen[v[1]] {
			varSeen[v[1]] = true
			promptVars = append(promptVars, v[1])
		}
	}

	framework := "vercel-ai-sdk"
	provider := "Cloud AI Provider"
	egressTarget := "external_cloud"

	if strings.Contains(content, "@ai-sdk/openai") || strings.Contains(content, "openai") {
		framework = "vercel-ai/openai"
		provider = "OpenAI"
	} else if strings.Contains(content, "@ai-sdk/anthropic") || strings.Contains(content, "@anthropic-ai/sdk") {
		framework = "vercel-ai/anthropic"
		provider = "Anthropic"
	} else if strings.Contains(content, "@langchain") {
		framework = "langchainjs"
	}

	var piiFindings []pii.Finding
	if systemPrompt != "" {
		piiFindings = append(piiFindings, pii.ScanPromptText(systemPrompt, 1)...)
	}
	piiFindings = append(piiFindings, pii.ScanVariables(promptVars, 1)...)

	if len(modelsFound) == 0 {
		modelsFound["vercel-ai-foundation-model"] = true
	}

	assetIdx := 0
	for modelName := range modelsFound {
		assetIdx++
		lowerModel := strings.ToLower(modelName)
		if strings.HasPrefix(lowerModel, "gpt") || strings.HasPrefix(lowerModel, "o1") || strings.HasPrefix(lowerModel, "o3") {
			provider = "OpenAI"
		} else if strings.HasPrefix(lowerModel, "claude") {
			provider = "Anthropic"
		}

		lineNum := 1
		for i, line := range lines {
			if strings.Contains(line, modelName) || tsAIImportRe.MatchString(line) {
				lineNum = i + 1
				break
			}
		}

		contextSnippet := ""
		start := lineNum - 2
		if start < 0 {
			start = 0
		}
		end := lineNum + 4
		if end > len(lines) {
			end = len(lines)
		}
		contextSnippet = strings.Join(lines[start:end], "\n")

		assetID := fmt.Sprintf("ts-%s-%d", filepath.Base(relPath), assetIdx)

		asset := DiscoveredAsset{
			ID:             assetID,
			Name:           fmt.Sprintf("%s (%s)", modelName, filepath.Base(relPath)),
			ModelFamily:    modelName,
			Provider:       provider,
			SourceFilePath: filepath.ToSlash(relPath),
			LineNumber:     lineNum,
			Language:       "typescript",
			Framework:      framework,
			ContextSnippet: contextSnippet,
			SystemPrompt:   systemPrompt,
			PromptVars:     promptVars,
			EgressTarget:   egressTarget,
			PIIFindings:    piiFindings,
			IntentKeywords: intentKeywords,
			FileHash:       fileHash,
		}

		assets = append(assets, asset)
	}

	return assets, nil
}
