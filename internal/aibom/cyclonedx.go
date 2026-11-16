package aibom

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/alexandrmotologa/governmesh/internal/scanner"
)

// CycloneDXProperty key-value attribute.
type CycloneDXProperty struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// CycloneDXTool represents tool metadata in CycloneDX 1.6.
type CycloneDXTool struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

// CycloneDXTools tools wrapper.
type CycloneDXTools struct {
	Components []CycloneDXTool `json:"components"`
}

// CycloneDXRootComponent represents the scanned repository.
type CycloneDXRootComponent struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

// CycloneDXMetadata header metadata.
type CycloneDXMetadata struct {
	Timestamp  string                 `json:"timestamp"`
	Tools      CycloneDXTools         `json:"tools"`
	Component  CycloneDXRootComponent `json:"component"`
	Properties []CycloneDXProperty    `json:"properties"`
}

// ModelParameters specifications.
type ModelParameters struct {
	Task               string `json:"task,omitempty"`
	ArchitectureFamily string `json:"architectureFamily,omitempty"`
}

// ModelConsiderations risk and use-case considerations.
type ModelConsiderations struct {
	Users     []string `json:"users,omitempty"`
	UseCases  []string `json:"useCases,omitempty"`
	RiskNotes string   `json:"riskNotes,omitempty"`
}

// CycloneDXModelCard ML model card schema.
type CycloneDXModelCard struct {
	ModelParameters     *ModelParameters     `json:"modelParameters,omitempty"`
	ModelConsiderations *ModelConsiderations `json:"modelConsiderations,omitempty"`
	Properties          []CycloneDXProperty  `json:"properties,omitempty"`
}

// CycloneDXComponent represents an AI component in the AIBOM.
type CycloneDXComponent struct {
	Type        string              `json:"type"` // "machine-learning-model", "framework"
	Name        string              `json:"name"`
	Version     string              `json:"version"`
	Description string              `json:"description,omitempty"`
	PURL        string              `json:"purl,omitempty"`
	ModelCard   *CycloneDXModelCard `json:"modelCard,omitempty"`
	Properties  []CycloneDXProperty `json:"properties,omitempty"`
}

// CycloneDXBOM 1.6 document.
type CycloneDXBOM struct {
	Schema       string               `json:"$schema"`
	BOMFormat    string               `json:"bomFormat"`
	SpecVersion  string               `json:"specVersion"`
	SerialNumber string               `json:"serialNumber"`
	Version      int                  `json:"version"`
	Metadata     CycloneDXMetadata    `json:"metadata"`
	Components   []CycloneDXComponent `json:"components"`
}

// GenerateCycloneDX synthesizes an authoritative CycloneDX 1.6 AIBOM document.
func GenerateCycloneDX(res *scanner.ScanResult) (*CycloneDXBOM, error) {
	uuidStr, err := generateUUIDv4()
	if err != nil {
		uuidStr = fmt.Sprintf("urn:uuid:gm-%d", time.Now().UnixNano())
	}

	metadataProps := []CycloneDXProperty{
		{Name: "governmesh:compliance:score", Value: fmt.Sprintf("%.1f", res.ComplianceScore)},
		{Name: "governmesh:compliance:passed_gate", Value: fmt.Sprintf("%t", res.PassedQualityGate)},
		{Name: "governmesh:attestation:hash_chain", Value: res.HashChain},
		{Name: "governmesh:risk:unacceptable", Value: fmt.Sprintf("%d", res.RiskSummary["unacceptable"])},
		{Name: "governmesh:risk:high", Value: fmt.Sprintf("%d", res.RiskSummary["high"])},
		{Name: "governmesh:risk:transparency", Value: fmt.Sprintf("%d", res.RiskSummary["transparency"])},
		{Name: "governmesh:risk:minimal", Value: fmt.Sprintf("%d", res.RiskSummary["minimal"])},
		{Name: "governmesh:pii:total_exposures", Value: fmt.Sprintf("%d", len(res.PIIExposures))},
	}

	var components []CycloneDXComponent
	for _, a := range res.Assets {
		sanitizedProvider := strings.ToLower(strings.ReplaceAll(a.Provider, " ", "-"))
		sanitizedModel := strings.ToLower(strings.ReplaceAll(a.ModelFamily, " ", "-"))
		purl := fmt.Sprintf("pkg:ml/%s/%s@latest", sanitizedProvider, sanitizedModel)

		cardProps := []CycloneDXProperty{
			{Name: "governmesh:risk:tier", Value: string(a.Risk.Tier)},
			{Name: "governmesh:risk:article", Value: a.Risk.ApplicableArticle},
			{Name: "governmesh:risk:justification", Value: a.Risk.Justification},
			{Name: "governmesh:egress:target", Value: a.EgressTarget},
			{Name: "governmesh:source:file", Value: a.SourceFilePath},
			{Name: "governmesh:source:line", Value: fmt.Sprintf("%d", a.LineNumber)},
			{Name: "governmesh:framework", Value: a.Framework},
		}

		if len(a.Risk.NISTCategories) > 0 {
			cardProps = append(cardProps, CycloneDXProperty{
				Name:  "governmesh:nist:categories",
				Value: strings.Join(a.Risk.NISTCategories, ","),
			})
		}
		if len(a.Risk.ISO42001Controls) > 0 {
			cardProps = append(cardProps, CycloneDXProperty{
				Name:  "governmesh:iso42001:controls",
				Value: strings.Join(a.Risk.ISO42001Controls, ","),
			})
		}

		comp := CycloneDXComponent{
			Type:        "machine-learning-model",
			Name:        a.ModelFamily,
			Version:     "1.0.0",
			Description: fmt.Sprintf("AI Model %s integrated via %s in %s", a.ModelFamily, a.Framework, a.SourceFilePath),
			PURL:        purl,
			ModelCard: &CycloneDXModelCard{
				ModelParameters: &ModelParameters{
					Task:               "text-and-multimodal-generation",
					ArchitectureFamily: "transformer",
				},
				ModelConsiderations: &ModelConsiderations{
					Users:     []string{"enterprise-application"},
					UseCases:  []string{a.Risk.ApplicableArticle},
					RiskNotes: a.Risk.Justification,
				},
				Properties: cardProps,
			},
			Properties: []CycloneDXProperty{
				{Name: "governmesh:file_hash", Value: a.FileHash},
			},
		}

		components = append(components, comp)
	}

	bom := &CycloneDXBOM{
		Schema:       "http://cyclonedx.org/schema/bom-1.6.schema.json",
		BOMFormat:    "CycloneDX",
		SpecVersion:  "1.6",
		SerialNumber: uuidStr,
		Version:      1,
		Metadata: CycloneDXMetadata{
			Timestamp: res.ScanTimestamp.UTC().Format(time.RFC3339),
			Tools: CycloneDXTools{
				Components: []CycloneDXTool{
					{
						Type:    "application",
						Name:    "governmesh",
						Version: "1.0.0",
					},
				},
			},
			Component: CycloneDXRootComponent{
				Type:    "application",
				Name:    res.ProjectName,
				Version: "1.0.0",
			},
			Properties: metadataProps,
		},
		Components: components,
	}

	return bom, nil
}

// ToJSON serializes CycloneDX BOM into indented JSON bytes.
func (b *CycloneDXBOM) ToJSON() ([]byte, error) {
	return json.MarshalIndent(b, "", "  ")
}

func generateUUIDv4() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	b[6] = (b[6] & 0x0f) | 0x40 // Version 4
	b[8] = (b[8] & 0x3f) | 0x80 // Variant RFC 4122
	return fmt.Sprintf("urn:uuid:%s-%s-%s-%s-%s",
		hex.EncodeToString(b[0:4]),
		hex.EncodeToString(b[4:6]),
		hex.EncodeToString(b[6:8]),
		hex.EncodeToString(b[8:10]),
		hex.EncodeToString(b[10:16])), nil
}
