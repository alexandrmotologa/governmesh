package scanner

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/alexandrmotologa/governmesh/internal/pii"
)

var knownGoAIPackages = map[string]string{
	"github.com/sashabaranov/go-openai":     "OpenAI",
	"github.com/tmc/langchaingo":             "LangChainGo",
	"google.golang.org/genai":               "Google Gemini",
	"github.com/anthropics/anthropic-sdk-go": "Anthropic",
	"github.com/cohere-ai/cohere-go":         "Cohere",
}

// ScanGoFile inspects a Go source file for AI client libraries using go/ast.
func ScanGoFile(path string, baseDir string) ([]DiscoveredAsset, error) {
	contentBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading go file: %w", err)
	}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, contentBytes, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parsing go AST: %w", err)
	}

	relPath, err := filepath.Rel(baseDir, path)
	if err != nil {
		relPath = path
	}

	hashSum := sha256.Sum256(contentBytes)
	fileHash := hex.EncodeToString(hashSum[:])

	var detectedProviders []string
	var detectedImports []string
	for _, imp := range node.Imports {
		importPath := strings.Trim(imp.Path.Value, `"`)
		for knownPkg, provider := range knownGoAIPackages {
			if strings.Contains(importPath, knownPkg) {
				detectedProviders = append(detectedProviders, provider)
				detectedImports = append(detectedImports, importPath)
			}
		}
	}

	if len(detectedProviders) == 0 {
		return nil, nil
	}

	lines := strings.Split(string(contentBytes), "\n")
	var intentKeywords []string
	modelsFound := make(map[string]bool)

	// Walk AST to find model strings and function names
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			if x.Name != nil {
				intentKeywords = append(intentKeywords, strings.ToLower(x.Name.Name))
			}
		case *ast.BasicLit:
			if x.Kind == token.STRING {
				val := strings.Trim(x.Value, `"`)
				lowerVal := strings.ToLower(val)
				if strings.HasPrefix(lowerVal, "gpt-") ||
					strings.HasPrefix(lowerVal, "claude-") ||
					strings.HasPrefix(lowerVal, "gemini-") ||
					strings.HasPrefix(lowerVal, "text-embedding-") {
					modelsFound[val] = true
				}
			}
		}
		return true
	})

	if len(modelsFound) == 0 {
		modelsFound["go-sdk-model"] = true
	}

	primaryProvider := detectedProviders[0]
	framework := "go-native-ai"
	if strings.Contains(strings.Join(detectedImports, " "), "langchaingo") {
		framework = "langchaingo"
	}

	var assets []DiscoveredAsset
	assetIdx := 0

	for modelName := range modelsFound {
		assetIdx++
		lineNum := 1
		for i, line := range lines {
			if strings.Contains(line, modelName) || strings.Contains(line, detectedImports[0]) {
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

		assetID := fmt.Sprintf("go-%s-%d", filepath.Base(relPath), assetIdx)

		// Scan intent keywords for PII
		var piiFindings []pii.Finding
		piiFindings = append(piiFindings, pii.ScanVariables(intentKeywords, lineNum)...)

		assets = append(assets, DiscoveredAsset{
			ID:             assetID,
			Name:           fmt.Sprintf("%s (%s)", modelName, filepath.Base(relPath)),
			ModelFamily:    modelName,
			Provider:       primaryProvider,
			SourceFilePath: filepath.ToSlash(relPath),
			LineNumber:     lineNum,
			Language:       "go",
			Framework:      framework,
			ContextSnippet: contextSnippet,
			EgressTarget:   "external_cloud",
			PIIFindings:    piiFindings,
			IntentKeywords: intentKeywords,
			FileHash:       fileHash,
		})
	}

	return assets, nil
}
