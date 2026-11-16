package report

import (
	"encoding/json"
	"fmt"

	"github.com/alexandrmotologa/governmesh/internal/config"
	"github.com/alexandrmotologa/governmesh/internal/scanner"
)

// SARIFReport represents OASIS SARIF v2.1.0 root structure.
type SARIFReport struct {
	Schema  string     `json:"$schema"`
	Version string     `json:"version"`
	Runs    []SARIFRun `json:"runs"`
}

type SARIFRun struct {
	Tool    SARIFTool     `json:"tool"`
	Results []SARIFResult `json:"results"`
}

type SARIFTool struct {
	Driver SARIFDriver `json:"driver"`
}

type SARIFDriver struct {
	Name           string      `json:"name"`
	Version        string      `json:"version"`
	InformationURI string      `json:"informationUri"`
	Rules          []SARIFRule `json:"rules"`
}

type SARIFRule struct {
	ID                   string                 `json:"id"`
	Name                 string                 `json:"name"`
	ShortDescription     SARIFText              `json:"shortDescription"`
	FullDescription      SARIFText              `json:"fullDescription,omitempty"`
	DefaultConfiguration SARIFRuleConfiguration `json:"defaultConfiguration"`
}

type SARIFRuleConfiguration struct {
	Level string `json:"level"` // "error", "warning", "note"
}

type SARIFText struct {
	Text string `json:"text"`
}

type SARIFResult struct {
	RuleID    string          `json:"ruleId"`
	Level     string          `json:"level"` // "error", "warning", "note"
	Message   SARIFText       `json:"message"`
	Locations []SARIFLocation `json:"locations"`
}

type SARIFLocation struct {
	PhysicalLocation SARIFPhysicalLocation `json:"physicalLocation"`
}

type SARIFPhysicalLocation struct {
	ArtifactLocation SARIFArtifactLocation `json:"artifactLocation"`
	Region           SARIFRegion           `json:"region"`
}

type SARIFArtifactLocation struct {
	URI   string `json:"uri"`
	Index int    `json:"index,omitempty"`
}

type SARIFRegion struct {
	StartLine   int `json:"startLine"`
	StartColumn int `json:"startColumn,omitempty"`
}

// GenerateSARIF synthesizes an OASIS SARIF 2.1.0 report for GitHub Security / PR annotations.
func GenerateSARIF(res *scanner.ScanResult) (*SARIFReport, error) {
	rulesMap := map[string]SARIFRule{
		"GM-EUAI-ART5": {
			ID:               "GM-EUAI-ART5",
			Name:             "EUAIActArticle5Prohibited",
			ShortDescription: SARIFText{Text: "EU AI Act Article 5 Prohibited AI Practice"},
			FullDescription:  SARIFText{Text: "The system implements prohibited capabilities (e.g. social scoring, subliminal manipulation, or unauthorized biometrics) banned under EU law."},
			DefaultConfiguration: SARIFRuleConfiguration{Level: "error"},
		},
		"GM-EUAI-ANNEX3": {
			ID:               "GM-EUAI-ANNEX3",
			Name:             "EUAIActAnnexIIIHighRisk",
			ShortDescription: SARIFText{Text: "EU AI Act Annex III High-Risk AI Classification"},
			FullDescription:  SARIFText{Text: "AI system deployed in sensitive domain (recruitment, credit scoring, healthcare, critical infrastructure) requiring strict conformity."},
			DefaultConfiguration: SARIFRuleConfiguration{Level: "warning"},
		},
		"GM-EUAI-ART50": {
			ID:               "GM-EUAI-ART50",
			Name:             "EUAIActArticle50Transparency",
			ShortDescription: SARIFText{Text: "EU AI Act Article 50 User Disclosure Required"},
			FullDescription:  SARIFText{Text: "Conversational bot or synthetic media generator requiring mandatory user disclosure notice."},
			DefaultConfiguration: SARIFRuleConfiguration{Level: "note"},
		},
		"GM-DATA-PII": {
			ID:               "GM-DATA-PII",
			Name:             "SensitiveDataEgressViolation",
			ShortDescription: SARIFText{Text: "Unmasked PII/PHI egressing to external cloud LLM"},
			FullDescription:  SARIFText{Text: "Personal identifier or health information variable bound into prompt template egressing to public cloud API."},
			DefaultConfiguration: SARIFRuleConfiguration{Level: "error"},
		},
	}

	var results []SARIFResult

	for _, a := range res.Assets {
		line := a.LineNumber
		if line < 1 {
			line = 1
		}

		// 1. EU AI Act findings
		var ruleID, level string
		switch a.Risk.Tier {
		case config.RiskUnacceptable:
			ruleID = "GM-EUAI-ART5"
			level = "error"
		case config.RiskHigh:
			ruleID = "GM-EUAI-ANNEX3"
			level = "warning"
		case config.RiskTransparency:
			ruleID = "GM-EUAI-ART50"
			level = "note"
		}

		if ruleID != "" {
			results = append(results, SARIFResult{
				RuleID: ruleID,
				Level:  level,
				Message: SARIFText{
					Text: fmt.Sprintf("[%s] %s: %s", a.Risk.Tier, a.ModelFamily, a.Risk.Justification),
				},
				Locations: []SARIFLocation{
					{
						PhysicalLocation: SARIFPhysicalLocation{
							ArtifactLocation: SARIFArtifactLocation{URI: a.SourceFilePath},
							Region:           SARIFRegion{StartLine: line, StartColumn: 1},
						},
					},
				},
			})
		}

		// 2. PII Findings
		for _, f := range a.PIIFindings {
			fLine := f.Line
			if fLine < 1 {
				fLine = line
			}
			piiLevel := "warning"
			if f.Severity == "critical" {
				piiLevel = "error"
			}
			results = append(results, SARIFResult{
				RuleID: "GM-DATA-PII",
				Level:  piiLevel,
				Message: SARIFText{
					Text: fmt.Sprintf("Sensitive data exposure (%s): %s in variable {%s}", f.Type, f.Description, f.Variable),
				},
				Locations: []SARIFLocation{
					{
						PhysicalLocation: SARIFPhysicalLocation{
							ArtifactLocation: SARIFArtifactLocation{URI: a.SourceFilePath},
							Region:           SARIFRegion{StartLine: fLine, StartColumn: 1},
						},
					},
				},
			})
		}
	}

	var rules []SARIFRule
	for _, r := range rulesMap {
		rules = append(rules, r)
	}

	report := &SARIFReport{
		Schema:  "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json",
		Version: "2.1.0",
		Runs: []SARIFRun{
			{
				Tool: SARIFTool{
					Driver: SARIFDriver{
						Name:           "governmesh",
						Version:        "1.0.0",
						InformationURI: "https://github.com/alexandrmotologa/governmesh",
						Rules:          rules,
					},
				},
				Results: results,
			},
		},
	}

	return report, nil
}

// ToJSON converts SARIF report into formatted JSON.
func (s *SARIFReport) ToJSON() ([]byte, error) {
	return json.MarshalIndent(s, "", "  ")
}
