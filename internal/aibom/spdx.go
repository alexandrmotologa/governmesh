package aibom

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/alexandrmotologa/governmesh/internal/scanner"
)

// SPDXElement represents an element in SPDX 3.0 document.
type SPDXElement struct {
	Type               string            `json:"type"` // "AIModel", "SoftwarePackage"
	SPDXID             string            `json:"spdxId"`
	Name               string            `json:"name"`
	Summary            string            `json:"summary,omitempty"`
	Description        string            `json:"description,omitempty"`
	Originator         string            `json:"originator,omitempty"`
	DownloadedFrom     string            `json:"downloadedFrom,omitempty"`
	Standard           string            `json:"standard,omitempty"`
	SafetyRiskCategory string            `json:"safetyRiskCategory,omitempty"`
	CustomProperties   map[string]string `json:"customProperties,omitempty"`
}

// SPDXRelationship links models and software modules.
type SPDXRelationship struct {
	From             string `json:"from"`
	To               string `json:"to"`
	RelationshipType string `json:"relationshipType"` // "DEPENDS_ON", "DESCRIBES"
}

// SPDXDocument represents an SPDX 3.0 AI profile manifest.
type SPDXDocument struct {
	SPDXVersion   string             `json:"spdxVersion"`
	DataLicense   string             `json:"dataLicense"`
	SPDXID        string             `json:"spdxId"`
	Name          string             `json:"name"`
	CreationInfo  map[string]string  `json:"creationInfo"`
	Elements      []SPDXElement      `json:"elements"`
	Relationships []SPDXRelationship `json:"relationships"`
}

// GenerateSPDX generates an SPDX 3.0 AI profile representation from scan results.
func GenerateSPDX(res *scanner.ScanResult) (*SPDXDocument, error) {
	rootID := "SPDXRef-RootPackage"
	elements := []SPDXElement{
		{
			Type:        "SoftwarePackage",
			SPDXID:      rootID,
			Name:        res.ProjectName,
			Description: fmt.Sprintf("Enterprise AI Application scanned from %s", res.ScannedPath),
			CustomProperties: map[string]string{
				"ComplianceScore": fmt.Sprintf("%.1f", res.ComplianceScore),
				"HashChain":       res.HashChain,
			},
		},
	}

	var relationships []SPDXRelationship

	for i, a := range res.Assets {
		elemID := fmt.Sprintf("SPDXRef-AIModel-%d", i+1)
		elem := SPDXElement{
			Type:               "AIModel",
			SPDXID:             elemID,
			Name:               a.ModelFamily,
			Summary:            fmt.Sprintf("AI integration via %s (%s)", a.Framework, a.Provider),
			Description:        a.Risk.Justification,
			Originator:         a.Provider,
			SafetyRiskCategory: string(a.Risk.Tier),
			CustomProperties: map[string]string{
				"EU_AI_Act_Article": a.Risk.ApplicableArticle,
				"SourceLocation":    fmt.Sprintf("%s:%d", a.SourceFilePath, a.LineNumber),
				"EgressTarget":      a.EgressTarget,
				"FileHash":          a.FileHash,
			},
		}
		elements = append(elements, elem)

		relationships = append(relationships, SPDXRelationship{
			From:             rootID,
			To:               elemID,
			RelationshipType: "DEPENDS_ON",
		})
	}

	doc := &SPDXDocument{
		SPDXVersion: "SPDX-3.0",
		DataLicense: "CC0-1.0",
		SPDXID:      "SPDXRef-DOCUMENT",
		Name:        fmt.Sprintf("AIBOM-%s", res.ProjectName),
		CreationInfo: map[string]string{
			"created":            time.Now().UTC().Format(time.RFC3339),
			"creators":           "Tool: governmesh-1.0.0",
			"licenseListVersion": "3.23",
		},
		Elements:      elements,
		Relationships: relationships,
	}

	return doc, nil
}

// ToJSON formats SPDX document as JSON bytes.
func (d *SPDXDocument) ToJSON() ([]byte, error) {
	return json.MarshalIndent(d, "", "  ")
}
