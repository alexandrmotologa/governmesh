package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// RiskLevel defines the string representation of EU AI Act risk categories.
type RiskLevel string

const (
	RiskUnacceptable RiskLevel = "unacceptable"
	RiskHigh         RiskLevel = "high"
	RiskTransparency RiskLevel = "transparency"
	RiskMinimal      RiskLevel = "minimal"
)

// Config represents enterprise governance rules and thresholds for governmesh.
type Config struct {
	Version      string           `yaml:"version" json:"version"`
	ProjectName  string           `yaml:"project_name" json:"project_name"`
	Scan         ScanConfig       `yaml:"scan" json:"scan"`
	Compliance   ComplianceConfig `yaml:"compliance" json:"compliance"`
	Attestation  AttestationConf  `yaml:"attestation" json:"attestation"`
	Exclusions   []string         `yaml:"exclusions" json:"exclusions"`
}

// ScanConfig defines directory scanning parameters.
type ScanConfig struct {
	IncludeExtensions []string `yaml:"include_extensions" json:"include_extensions"`
	ExcludePaths      []string `yaml:"exclude_paths" json:"exclude_paths"`
	MaxFileSizeMB     int64    `yaml:"max_file_size_mb" json:"max_file_size_mb"`
	ScanPII           bool     `yaml:"scan_pii" json:"scan_pii"`
	ScanManifests     bool     `yaml:"scan_manifests" json:"scan_manifests"`
}

// ComplianceConfig defines risk thresholds and CI/CD quality gate enforcement.
type ComplianceConfig struct {
	FailOnRiskLevel    RiskLevel `yaml:"fail_on_risk_level" json:"fail_on_risk_level"`
	FailOnPIIEgress    bool      `yaml:"fail_on_pii_egress" json:"fail_on_pii_egress"`
	EnforceHumanReview bool      `yaml:"enforce_human_review" json:"enforce_human_review"`
	TargetStandard     string    `yaml:"target_standard" json:"target_standard"` // "eu_ai_act", "nist_rmf", "iso_42001"
}

// AttestationConf defines cryptographic signing preferences.
type AttestationConf struct {
	SignReport   bool   `yaml:"sign_report" json:"sign_report"`
	SignerName   string `yaml:"signer_name" json:"signer_name"`
	Organization string `yaml:"organization" json:"organization"`
}

// DefaultConfig returns safe, standard defaults for enterprise scanning.
func DefaultConfig() *Config {
	return &Config{
		Version:     "1.0.0",
		ProjectName: "Default Organization",
		Scan: ScanConfig{
			IncludeExtensions: []string{".py", ".ts", ".tsx", ".js", ".jsx", ".go", ".yaml", ".yml", ".json", ".env", "Dockerfile"},
			ExcludePaths: []string{
				".git", "node_modules", "vendor", "__pycache__", ".venv", "dist", "build", ".next",
			},
			MaxFileSizeMB: 10,
			ScanPII:       true,
			ScanManifests: true,
		},
		Compliance: ComplianceConfig{
			FailOnRiskLevel:    RiskUnacceptable,
			FailOnPIIEgress:    true,
			EnforceHumanReview: true,
			TargetStandard:     "eu_ai_act",
		},
		Attestation: AttestationConf{
			SignReport:   true,
			SignerName:   "governmesh-scanner",
			Organization: "Enterprise AI Security Office",
		},
		Exclusions: []string{},
	}
}

// LoadConfig loads configuration from a YAML file path or returns defaults if path is empty.
func LoadConfig(path string) (*Config, error) {
	cfg := DefaultConfig()
	if path == "" {
		return cfg, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing yaml configuration: %w", err)
	}

	return cfg, nil
}
