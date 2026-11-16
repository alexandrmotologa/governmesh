package pii

// Finding represents sensitive data detection inside a prompt or variable.
type Finding struct {
	Type        string `json:"type"`        // "ssn", "iban", "email", "phone", "credit_card", "phi", "biometric", "financial"
	Snippet     string `json:"snippet"`     // Redacted snippet
	Line        int    `json:"line"`
	Severity    string `json:"severity"`    // "critical", "high", "medium", "low"
	Variable    string `json:"variable"`    // Variable name if identified
	Description string `json:"description"`
}

// Exposure aggregates sensitive variables egressing to external cloud providers.
type Exposure struct {
	AssetID      string `json:"asset_id"`
	FilePath     string `json:"file_path"`
	Line         int    `json:"line"`
	PIIType      string `json:"pii_type"`
	TargetCloud  string `json:"target_cloud"`
	Snippet      string `json:"snippet"`
	IsProhibited bool   `json:"is_prohibited"`
}
