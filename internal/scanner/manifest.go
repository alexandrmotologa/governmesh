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
	baseImageRe   = regexp.MustCompile(`(?i)^\s*FROM\s+([a-zA-Z0-9_\-\.\/:]+)`)
	weightsDownloadRe = regexp.MustCompile(`(?i)(?:huggingface-cli\s+download|git\s+lfs\s+pull|curl.*?\/models\/|wget.*?\/models\/|[a-zA-Z0-9_\-\.\/]+\.(?:safetensors|onnx|gguf|bin|pt))`)
	envModelRe     = regexp.MustCompile(`(?i)(?:MODEL_NAME|DEFAULT_MODEL|LLM_MODEL)\s*=\s*["']?([a-zA-Z0-9_\-\.\/]+)`)
	envSecretRe    = regexp.MustCompile(`(?i)(?:OPENAI_API_KEY|ANTHROPIC_API_KEY|COHERE_API_KEY|HUGGING_FACE_HUB_TOKEN|HF_TOKEN)\s*=\s*["']?([^\s"']+)`)
	k8sImageRe     = regexp.MustCompile(`image:\s*["']?([a-zA-Z0-9_\-\.\/:]+)`)
)

// ScanManifestFile inspects Dockerfile, Compose, Helm, or .env files for infrastructure AI signals.
func ScanManifestFile(path string, baseDir string) ([]DiscoveredAsset, error) {
	contentBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading manifest file: %w", err)
	}

	content := string(contentBytes)
	fileName := filepath.Base(path)
	lowerName := strings.ToLower(fileName)

	isDocker := strings.Contains(lowerName, "dockerfile")
	isCompose := strings.Contains(lowerName, "compose") || strings.Contains(lowerName, "docker-compose")
	isEnv := strings.HasPrefix(lowerName, ".env")
	isYaml := strings.HasSuffix(lowerName, ".yaml") || strings.HasSuffix(lowerName, ".yml")

	if !isDocker && !isCompose && !isEnv && !isYaml {
		return nil, nil
	}

	relPath, err := filepath.Rel(baseDir, path)
	if err != nil {
		relPath = path
	}

	hashSum := sha256.Sum256(contentBytes)
	fileHash := hex.EncodeToString(hashSum[:])

	var assets []DiscoveredAsset

	// 1. Dockerfile inspections
	if isDocker {
		lines := strings.Split(content, "\n")
		for i, line := range lines {
			if m := baseImageRe.FindStringSubmatch(line); len(m) > 1 {
				img := m[1]
				lowerImg := strings.ToLower(img)
				if strings.Contains(lowerImg, "ollama") ||
					strings.Contains(lowerImg, "vllm") ||
					strings.Contains(lowerImg, "text-generation-inference") ||
					strings.Contains(lowerImg, "tritonserver") ||
					strings.Contains(lowerImg, "pytorch") {

					provider := "Self-Hosted / Container Runtime"
					framework := "container-runtime"
					egress := "local_runtime"

					if strings.Contains(lowerImg, "ollama") {
						framework = "ollama-container"
					} else if strings.Contains(lowerImg, "vllm") {
						framework = "vllm-container"
					}

					asset := DiscoveredAsset{
						ID:             fmt.Sprintf("dock-%s-%d", filepath.Base(relPath), i+1),
						Name:           fmt.Sprintf("Container Model Runtime (%s)", img),
						ModelFamily:    img,
						Provider:       provider,
						SourceFilePath: filepath.ToSlash(relPath),
						LineNumber:     i + 1,
						Language:       "manifest",
						Framework:      framework,
						ContextSnippet: line,
						EgressTarget:   egress,
						FileHash:       fileHash,
					}
					assets = append(assets, asset)
				}
			}
		}
	}

	// 2. Model weights detection in build manifests
	if weightMatches := weightsDownloadRe.FindAllString(content, -1); len(weightMatches) > 0 {
		weightFile := weightMatches[0]
		asset := DiscoveredAsset{
			ID:             fmt.Sprintf("weight-%s", filepath.Base(relPath)),
			Name:           fmt.Sprintf("Embedded Model Weights (%s)", filepath.Base(weightFile)),
			ModelFamily:    weightFile,
			Provider:       "Model Checkpoint / Weights",
			SourceFilePath: filepath.ToSlash(relPath),
			LineNumber:     1,
			Language:       "manifest",
			Framework:      "weights-artifact",
			ContextSnippet: weightFile,
			EgressTarget:   "local_runtime",
			FileHash:       fileHash,
		}
		assets = append(assets, asset)
	}

	// 3. Environment variables with model names or API keys
	if modelMatches := envModelRe.FindAllStringSubmatch(content, -1); len(modelMatches) > 0 {
		for _, m := range modelMatches {
			if len(m) > 1 {
				modelName := m[1]
				asset := DiscoveredAsset{
					ID:             fmt.Sprintf("env-model-%s", modelName),
					Name:           fmt.Sprintf("Env Configured Model: %s", modelName),
					ModelFamily:    modelName,
					Provider:       "Environment Injected Provider",
					SourceFilePath: filepath.ToSlash(relPath),
					LineNumber:     1,
					Language:       "manifest",
					Framework:      "env-config",
					ContextSnippet: m[0],
					EgressTarget:   "external_cloud",
					FileHash:       fileHash,
				}
				assets = append(assets, asset)
			}
		}
	}

	// 4. Secret keys leak inspection
	if secretMatches := envSecretRe.FindAllStringSubmatch(content, -1); len(secretMatches) > 0 {
		for _, m := range secretMatches {
			if len(m) > 1 && !strings.Contains(m[1], "${") && len(m[1]) > 10 {
				finding := pii.ScanPromptText(m[0], 1)
				if len(finding) > 0 {
					asset := DiscoveredAsset{
						ID:             fmt.Sprintf("env-secret-%s", filepath.Base(relPath)),
						Name:           "Plaintext Cloud AI Secret in Configuration",
						ModelFamily:    "Cloud API Credential",
						Provider:       "External AI Gateway",
						SourceFilePath: filepath.ToSlash(relPath),
						LineNumber:     1,
						Language:       "manifest",
						Framework:      "secret-exposure",
						ContextSnippet: "SECRET_KEY=[REDACTED]",
						PIIFindings:    finding,
						FileHash:       fileHash,
					}
					assets = append(assets, asset)
				}
			}
		}
	}

	return assets, nil
}
