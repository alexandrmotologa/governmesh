package aibom

import (
	"encoding/json"
	"fmt"
	"os"
)

// RiskDiff aggregates differences between two CycloneDX AIBOM snapshots.
type RiskDiff struct {
	BaseProject   string   `json:"base_project"`
	TargetProject string   `json:"target_project"`
	AddedModels   []string `json:"added_models"`
	RemovedModels []string `json:"removed_models"`
	RiskEscalations []string `json:"risk_escalations"`
	BaseScore     string   `json:"base_score"`
	TargetScore   string   `json:"target_score"`
	ScoreDelta    float64  `json:"score_delta"`
}

// CompareAIBOMFiles compares two CycloneDX 1.6 JSON files and detects risk regressions.
func CompareAIBOMFiles(basePath, targetPath string) (*RiskDiff, error) {
	baseData, err := os.ReadFile(basePath)
	if err != nil {
		return nil, fmt.Errorf("reading base aibom: %w", err)
	}
	targetData, err := os.ReadFile(targetPath)
	if err != nil {
		return nil, fmt.Errorf("reading target aibom: %w", err)
	}

	var baseBOM, targetBOM CycloneDXBOM
	if err := json.Unmarshal(baseData, &baseBOM); err != nil {
		return nil, fmt.Errorf("parsing base aibom json: %w", err)
	}
	if err := json.Unmarshal(targetData, &targetBOM); err != nil {
		return nil, fmt.Errorf("parsing target aibom json: %w", err)
	}

	diff := &RiskDiff{
		BaseProject:   baseBOM.Metadata.Component.Name,
		TargetProject: targetBOM.Metadata.Component.Name,
	}

	baseModels := make(map[string]string)
	for _, c := range baseBOM.Components {
		tier := "minimal"
		if c.ModelCard != nil {
			for _, p := range c.ModelCard.Properties {
				if p.Name == "governmesh:risk:tier" {
					tier = p.Value
				}
			}
		}
		baseModels[c.Name] = tier
	}

	targetModels := make(map[string]string)
	for _, c := range targetBOM.Components {
		tier := "minimal"
		if c.ModelCard != nil {
			for _, p := range c.ModelCard.Properties {
				if p.Name == "governmesh:risk:tier" {
					tier = p.Value
				}
			}
		}
		targetModels[c.Name] = tier
		if _, existed := baseModels[c.Name]; !existed {
			diff.AddedModels = append(diff.AddedModels, fmt.Sprintf("%s (Risk: %s)", c.Name, tier))
			if tier == "high" || tier == "unacceptable" {
				diff.RiskEscalations = append(diff.RiskEscalations, fmt.Sprintf("Newly added model '%s' has %s risk level", c.Name, tier))
			}
		}
	}

	for name, tier := range baseModels {
		if _, exists := targetModels[name]; !exists {
			diff.RemovedModels = append(diff.RemovedModels, fmt.Sprintf("%s (was: %s)", name, tier))
		}
	}

	return diff, nil
}
