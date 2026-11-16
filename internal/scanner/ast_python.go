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
	pyAIImportRe = regexp.MustCompile(`(?m)^\s*(?:import\s+(openai|anthropic|langchain|transformers|litellm|google\.generativeai|cohere|mistralai|vllm|ollama)|from\s+(openai|anthropic|langchain[a-zA-Z0-9_\.]*|transformers|litellm|google|cohere|mistralai|vllm|ollama)\s+import)`)
	pyModelArgRe = regexp.MustCompile(`(?:model|model_name)\s*=\s*["']([a-zA-Z0-9\._\-\/]+)["']`)
	pySystemPromptRe = regexp.MustCompile(`(?s)(?:system_prompt|system_message|SYSTEM_PROMPT)\s*=\s*["']{1,3}(.*?)["']{1,3}`)
	pyRoleSystemRe = regexp.MustCompile(`(?s)["']role["']\s*:\s*["']system["']\s*,\s*["']content["']\s*:\s*["']{1,3}(.*?)["']{1,3}`)
	pyTemplateVarRe = regexp.MustCompile(`\{([a-zA-Z0-9_]+)\}`)
	pyCallRe = regexp.MustCompile(`([a-zA-Z0-9_]+)\.(?:chat\.completions\.create|messages\.create|generate|invoke|predict|call)\s*\(`)
	pyFunctionDefRe = regexp.MustCompile(`(?m)^\s*def\s+([a-zA-Z0-9_]+)\s*\(`)
)

// ScanPythonFile inspects a Python file for AI SDK usage, model references, and prompts.
func ScanPythonFile(path string, baseDir string) ([]DiscoveredAsset, error) {
	contentBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading python file: %w", err)
	}

	content := string(contentBytes)
	hasAIImport := pyAIImportRe.MatchString(content)
	hasModelRef := pyModelArgRe.MatchString(content)

	// If no AI imports or model references, skip this file
	if !hasAIImport && !hasModelRef {
		return nil, nil
	}

	relPath, err := filepath.Rel(baseDir, path)
	if err != nil {
		relPath = path
	}

	// Compute file hash for cryptographic attestation
	hashSum := sha256.Sum256(contentBytes)
	fileHash := hex.EncodeToString(hashSum[:])

	lines := strings.Split(content, "\n")
	var assets []DiscoveredAsset

	// Extract function names for domain intent detection
	var intentKeywords []string
	fnMatches := pyFunctionDefRe.FindAllStringSubmatch(content, -1)
	for _, m := range fnMatches {
		if len(m) > 1 {
			intentKeywords = append(intentKeywords, strings.ToLower(m[1]))
		}
	}

	// Extract models referenced
	modelMatches := pyModelArgRe.FindAllStringSubmatch(content, -1)
	modelsFound := make(map[string]bool)
	for _, m := range modelMatches {
		if len(m) > 1 {
			modelsFound[m[1]] = true
		}
	}

	// Extract system prompts
	systemPrompt := ""
	if sm := pySystemPromptRe.FindStringSubmatch(content); len(sm) > 1 {
		systemPrompt = strings.TrimSpace(sm[1])
	} else if sm2 := pyRoleSystemRe.FindStringSubmatch(content); len(sm2) > 1 {
		systemPrompt = strings.TrimSpace(sm2[1])
	}

	// Extract prompt template variables
	var promptVars []string
	varSeen := make(map[string]bool)
	varMatches := pyTemplateVarRe.FindAllStringSubmatch(content, -1)
	for _, v := range varMatches {
		if len(v) > 1 && !varSeen[v[1]] {
			varSeen[v[1]] = true
			promptVars = append(promptVars, v[1])
		}
	}

	// Detect Framework & Provider
	framework := "generic-ai"
	provider := "Unknown AI Provider"
	egressTarget := "external_cloud"

	if strings.Contains(content, "langchain") {
		framework = "langchain"
	} else if strings.Contains(content, "openai") {
		framework = "openai-sdk"
		provider = "OpenAI"
	} else if strings.Contains(content, "anthropic") {
		framework = "anthropic-sdk"
		provider = "Anthropic"
	} else if strings.Contains(content, "transformers") {
		framework = "huggingface-transformers"
		provider = "Hugging Face (Local/Self-Hosted)"
		egressTarget = "local_runtime"
	} else if strings.Contains(content, "ollama") {
		framework = "ollama"
		provider = "Ollama (Self-Hosted)"
		egressTarget = "local_runtime"
	} else if strings.Contains(content, "vllm") {
		framework = "vllm"
		provider = "vLLM (Self-Hosted)"
		egressTarget = "local_runtime"
	}

	// Scan PII in system prompt and variables
	var piiFindings []pii.Finding
	if systemPrompt != "" {
		piiFindings = append(piiFindings, pii.ScanPromptText(systemPrompt, 1)...)
	}
	piiFindings = append(piiFindings, pii.ScanVariables(promptVars, 1)...)

	// Build asset representation for each distinct model found (or general file asset)
	if len(modelsFound) == 0 {
		modelsFound["default-foundation-model"] = true
	}

	assetIdx := 0
	for modelName := range modelsFound {
		assetIdx++
		// If model name reveals provider, refine it
		lowerModel := strings.ToLower(modelName)
		if strings.HasPrefix(lowerModel, "gpt") || strings.HasPrefix(lowerModel, "o1") || strings.HasPrefix(lowerModel, "o3") || strings.Contains(lowerModel, "text-embedding") {
			provider = "OpenAI"
			egressTarget = "external_cloud"
		} else if strings.HasPrefix(lowerModel, "claude") {
			provider = "Anthropic"
			egressTarget = "external_cloud"
		} else if strings.HasPrefix(lowerModel, "gemini") {
			provider = "Google Cloud"
			egressTarget = "external_cloud"
		} else if strings.HasPrefix(lowerModel, "llama") || strings.HasPrefix(lowerModel, "mistral") || strings.HasPrefix(lowerModel, "qwen") || strings.HasPrefix(lowerModel, "deepseek") {
			if egressTarget != "external_cloud" {
				provider = "Open-Source Weights"
			}
		}

		lineNum := 1
		for i, line := range lines {
			if strings.Contains(line, modelName) || pyAIImportRe.MatchString(line) {
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

		assetID := fmt.Sprintf("py-%s-%d", filepath.Base(relPath), assetIdx)

		asset := DiscoveredAsset{
			ID:             assetID,
			Name:           fmt.Sprintf("%s (%s)", modelName, filepath.Base(relPath)),
			ModelFamily:    modelName,
			Provider:       provider,
			SourceFilePath: filepath.ToSlash(relPath),
			LineNumber:     lineNum,
			Language:       "python",
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
