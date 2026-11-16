package scanner

import (
	"time"

	"github.com/alexandrmotologa/governmesh/internal/config"
	"github.com/alexandrmotologa/governmesh/internal/pii"
)

// RiskAssessment holds compliance evaluations across regulations.
type RiskAssessment struct {
	Tier              config.RiskLevel `json:"tier"` // unacceptable, high, transparency, minimal
	ApplicableArticle string           `json:"applicable_article"`
	Justification     string           `json:"justification"`
	ProhibitedReason  string           `json:"prohibited_reason,omitempty"`
	RequiresHumanGate bool             `json:"requires_human_gate"`
	NISTCategories    []string         `json:"nist_categories"`
	ISO42001Controls  []string         `json:"iso_42001_controls"`
}

// DiscoveredAsset represents an identified AI system, model call, or component.
type DiscoveredAsset struct {
	ID             string                 `json:"id"`
	Name           string                 `json:"name"`
	ModelFamily    string                 `json:"model_family"`
	Provider       string                 `json:"provider"`
	SourceFilePath string                 `json:"source_file_path"`
	LineNumber     int                    `json:"line_number"`
	Language       string                 `json:"language"`  // "python", "typescript", "go", "manifest", "prompt"
	Framework      string                 `json:"framework"` // "langchain", "openai-sdk", "anthropic-sdk", "vercel-ai", etc.
	ContextSnippet string                 `json:"context_snippet"`
	Parameters     map[string]interface{} `json:"parameters,omitempty"`
	SystemPrompt   string                 `json:"system_prompt,omitempty"`
	PromptVars     []string               `json:"prompt_vars,omitempty"`
	EgressTarget   string                 `json:"egress_target,omitempty"` // "external_cloud", "local_runtime", "private_vpc"
	PIIFindings    []pii.Finding          `json:"pii_findings,omitempty"`
	IntentKeywords []string               `json:"intent_keywords,omitempty"`
	Risk           RiskAssessment         `json:"risk"`
	FileHash       string                 `json:"file_hash"`
}

// ScanResult contains the full scan inventory, risk classifications, and audit status.
type ScanResult struct {
	ProjectName       string                   `json:"project_name"`
	ScannedPath       string                   `json:"scanned_path"`
	ScanTimestamp     time.Time                `json:"scan_timestamp"`
	DurationMs        int64                    `json:"duration_ms"`
	TotalFiles        int                      `json:"total_files"`
	ScannedAIFiles    int                      `json:"scanned_ai_files"`
	Assets            []DiscoveredAsset        `json:"assets"`
	RiskSummary       map[config.RiskLevel]int `json:"risk_summary"`
	PIIExposures      []pii.Exposure           `json:"pii_exposures"`
	ComplianceScore   float64                  `json:"compliance_score"` // 0 - 100
	HashChain         string                   `json:"hash_chain"`       // Merkle root of all scanned components
	PassedQualityGate bool                     `json:"passed_quality_gate"`
	GateFailures      []string                 `json:"gate_failures,omitempty"`
}
