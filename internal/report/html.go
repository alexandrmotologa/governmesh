package report

import (
	"fmt"
	"html"
	"strings"
	"time"

	"github.com/alexandrmotologa/governmesh/internal/attestation"
	"github.com/alexandrmotologa/governmesh/internal/config"
	"github.com/alexandrmotologa/governmesh/internal/scanner"
)

// GenerateHTMLReport creates a standalone, self-contained HTML audit document.
func GenerateHTMLReport(res *scanner.ScanResult, att *attestation.Manifest) string {
	gateBadgeClass := "badge-pass"
	gateStatusText := "PASSED"
	if !res.PassedQualityGate {
		gateBadgeClass = "badge-fail"
		gateStatusText = "NON-COMPLIANT"
	}

	var sb strings.Builder
	sb.WriteString(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>AI Governance & Regulatory Audit Report - ` + html.EscapeString(res.ProjectName) + `</title>
<style>
  :root {
    --bg: #090d16;
    --surface: #111827;
    --surface-border: #1f2937;
    --text: #f3f4f6;
    --text-muted: #9ca3af;
    --primary: #3b82f6;
    --danger: #ef4444;
    --warning: #f59e0b;
    --success: #10b981;
    --purple: #8b5cf6;
  }
  * { box-sizing: border-box; margin: 0; padding: 0; }
  body {
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
    background-color: var(--bg);
    color: var(--text);
    padding: 2.5rem 1.5rem;
    line-height: 1.5;
  }
  .container { max-width: 1080px; margin: 0 auto; }
  header {
    border-bottom: 1px solid var(--surface-border);
    padding-bottom: 1.5rem;
    margin-bottom: 2rem;
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
  }
  h1 { font-size: 1.85rem; font-weight: 700; color: #fff; margin-bottom: 0.35rem; }
  .subtitle { color: var(--text-muted); font-size: 0.95rem; }
  .badge {
    padding: 0.35rem 0.85rem;
    border-radius: 9999px;
    font-size: 0.85rem;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.05em;
  }
  .badge-pass { background-color: rgba(16, 185, 129, 0.2); color: var(--success); border: 1px solid var(--success); }
  .badge-fail { background-color: rgba(239, 68, 68, 0.2); color: var(--danger); border: 1px solid var(--danger); }
  .grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(220px, 1fr)); gap: 1rem; margin-bottom: 2rem; }
  .card {
    background: var(--surface);
    border: 1px solid var(--surface-border);
    border-radius: 8px;
    padding: 1.25rem;
  }
  .card-title { font-size: 0.8rem; text-transform: uppercase; color: var(--text-muted); font-weight: 600; margin-bottom: 0.5rem; }
  .card-value { font-size: 1.75rem; font-weight: 700; color: #fff; }
  section { margin-bottom: 2.5rem; }
  h2 { font-size: 1.25rem; font-weight: 600; color: #fff; margin-bottom: 1rem; border-left: 3px solid var(--primary); padding-left: 0.75rem; }
  table { width: 100%; border-collapse: collapse; background: var(--surface); border: 1px solid var(--surface-border); border-radius: 8px; overflow: hidden; margin-bottom: 1rem; }
  th, td { padding: 0.75rem 1rem; text-align: left; font-size: 0.875rem; }
  th { background: #1a2234; color: var(--text-muted); font-weight: 600; border-bottom: 1px solid var(--surface-border); }
  tr:not(:last-child) td { border-bottom: 1px solid var(--surface-border); }
  code { font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace; font-size: 0.825rem; background: rgba(255,255,255,0.06); padding: 0.15rem 0.4rem; border-radius: 4px; }
  .tier-unacceptable { color: #f87171; font-weight: 600; }
  .tier-high { color: #fbbf24; font-weight: 600; }
  .tier-transparency { color: #60a5fa; font-weight: 600; }
  .tier-minimal { color: #34d399; font-weight: 600; }
  .attestation-block { background: #0c111d; border: 1px solid var(--surface-border); border-radius: 8px; padding: 1.25rem; font-family: monospace; font-size: 0.8rem; overflow-x: auto; color: #94a3b8; }
  @media print {
    body { background: #fff; color: #111; padding: 0; }
    .card, table, .attestation-block { border-color: #ddd; background: #fff; color: #111; }
    th { background: #f3f4f6; color: #333; }
    code { background: #eee; color: #000; }
  }
</style>
</head>
<body>
<div class="container">
  <header>
    <div>
      <h1>AI Regulatory Compliance & Governance Audit</h1>
      <div class="subtitle">Target: <strong>` + html.EscapeString(res.ProjectName) + `</strong> | Date: ` + res.ScanTimestamp.Format(time.RFC1123) + `</div>
    </div>
    <div>
      <span class="badge ` + gateBadgeClass + `">` + gateStatusText + `</span>
    </div>
  </header>

  <div class="grid">
    <div class="card">
      <div class="card-title">Compliance Score</div>
      <div class="card-value" style="color: ` + scoreColor(res.ComplianceScore) + `">` + fmt.Sprintf("%.1f%%", res.ComplianceScore) + `</div>
    </div>
    <div class="card">
      <div class="card-title">AI Models Discovered</div>
      <div class="card-value">` + fmt.Sprintf("%d", len(res.Assets)) + `</div>
    </div>
    <div class="card">
      <div class="card-title">High & Unacceptable Risks</div>
      <div class="card-value" style="color: var(--danger)">` + fmt.Sprintf("%d", res.RiskSummary[config.RiskUnacceptable]+res.RiskSummary[config.RiskHigh]) + `</div>
    </div>
    <div class="card">
      <div class="card-title">PII Data Exposures</div>
      <div class="card-value" style="color: var(--warning)">` + fmt.Sprintf("%d", len(res.PIIExposures)) + `</div>
    </div>
  </div>
`)

	if len(res.GateFailures) > 0 {
		sb.WriteString(`
  <section>
    <h2>Compliance Gate Warnings</h2>
    <div style="background: rgba(239,68,68,0.1); border: 1px solid var(--danger); border-radius: 8px; padding: 1rem;">
      <ul style="padding-left: 1.5rem; color: #fca5a5;">
`)
		for _, f := range res.GateFailures {
			sb.WriteString(`<li>` + html.EscapeString(f) + `</li>`)
		}
		sb.WriteString(`
      </ul>
    </div>
  </section>
`)
	}

	// EU AI Act Breakdown
	sb.WriteString(`
  <section>
    <h2>EU AI Act Risk Distribution</h2>
    <table>
      <thead>
        <tr>
          <th>Risk Tier</th>
          <th>Regulatory Basis</th>
          <th>Discovered Systems</th>
          <th>Conformity Mandate</th>
        </tr>
      </thead>
      <tbody>
        <tr>
          <td><span class="tier-unacceptable">Unacceptable Risk</span></td>
          <td>Article 5 (Prohibited AI Practices)</td>
          <td><strong>` + fmt.Sprintf("%d", res.RiskSummary[config.RiskUnacceptable]) + `</strong></td>
          <td>Hard CI/CD gate block & prohibited use remediation</td>
        </tr>
        <tr>
          <td><span class="tier-high">High-Risk</span></td>
          <td>Article 6 & Annex III (Employment, Credit, Critical)</td>
          <td><strong>` + fmt.Sprintf("%d", res.RiskSummary[config.RiskHigh]) + `</strong></td>
          <td>Risk management, data governance, human oversight</td>
        </tr>
        <tr>
          <td><span class="tier-transparency">Specific Transparency</span></td>
          <td>Article 50 (Chatbots, Synthetic Media)</td>
          <td><strong>` + fmt.Sprintf("%d", res.RiskSummary[config.RiskTransparency]) + `</strong></td>
          <td>Mandatory AI user disclosure notices</td>
        </tr>
        <tr>
          <td><span class="tier-minimal">Minimal Risk</span></td>
          <td>Article 69 (General Internal Use)</td>
          <td><strong>` + fmt.Sprintf("%d", res.RiskSummary[config.RiskMinimal]) + `</strong></td>
          <td>Permitted with standard AIBOM tracking</td>
        </tr>
      </tbody>
    </table>
  </section>
`)

	// Discovered Inventory
	sb.WriteString(`
  <section>
    <h2>Discovered AI Bill of Materials (AIBOM)</h2>
    <table>
      <thead>
        <tr>
          <th>Model / Component</th>
          <th>Provider</th>
          <th>Framework</th>
          <th>Risk Tier</th>
          <th>Egress Target</th>
          <th>Source Code Location</th>
        </tr>
      </thead>
      <tbody>
`)
	for _, a := range res.Assets {
		tierClass := "tier-" + string(a.Risk.Tier)
		sb.WriteString(fmt.Sprintf(`
        <tr>
          <td><strong>%s</strong></td>
          <td>%s</td>
          <td><code>%s</code></td>
          <td><span class="%s">%s</span></td>
          <td><code>%s</code></td>
          <td><code>%s:%d</code></td>
        </tr>`,
			html.EscapeString(a.ModelFamily),
			html.EscapeString(a.Provider),
			html.EscapeString(a.Framework),
			tierClass,
			html.EscapeString(string(a.Risk.Tier)),
			html.EscapeString(a.EgressTarget),
			html.EscapeString(a.SourceFilePath),
			a.LineNumber,
		))
	}
	sb.WriteString(`
      </tbody>
    </table>
  </section>
`)

	// Sensitive Data Exposures
	sb.WriteString(`
  <section>
    <h2>Sensitive Data & PII Egress Audit</h2>
`)
	if len(res.PIIExposures) == 0 {
		sb.WriteString(`<p style="color: var(--text-muted);">No sensitive personal data variables detected egressing to external cloud AI models.</p>`)
	} else {
		sb.WriteString(`
    <table>
      <thead>
        <tr>
          <th>Source File</th>
          <th>Line</th>
          <th>PII Type</th>
          <th>Target Cloud</th>
          <th>Redacted Sample</th>
          <th>Prohibited Status</th>
        </tr>
      </thead>
      <tbody>
`)
		for _, p := range res.PIIExposures {
			prohLabel := `<span style="color: var(--text-muted)">Permitted</span>`
			if p.IsProhibited {
				prohLabel = `<span style="color: var(--danger); font-weight: bold;">VIOLATION</span>`
			}
			sb.WriteString(fmt.Sprintf(`
        <tr>
          <td><code>%s</code></td>
          <td>%d</td>
          <td><code>%s</code></td>
          <td>%s</td>
          <td><code>%s</code></td>
          <td>%s</td>
        </tr>`,
				html.EscapeString(p.FilePath),
				p.Line,
				html.EscapeString(p.PIIType),
				html.EscapeString(p.TargetCloud),
				html.EscapeString(p.Snippet),
				prohLabel,
			))
		}
		sb.WriteString(`
      </tbody>
    </table>
`)
	}
	sb.WriteString(`
  </section>
`)

	// Cryptographic Attestation Block
	if att != nil {
		sb.WriteString(`
  <section>
    <h2>Cryptographic Attestation & Audit Trail</h2>
    <div class="attestation-block">
      <div>Algorithm: ` + html.EscapeString(att.Algorithm) + `</div>
      <div>Signer: ` + html.EscapeString(att.Signer) + ` (` + html.EscapeString(att.Organization) + `)</div>
      <div>Attestation Timestamp: ` + html.EscapeString(att.Timestamp) + `</div>
      <div>Merkle Root Hash: <code>` + html.EscapeString(att.RootHash) + `</code></div>
      <div>AIBOM Checksum: <code>` + html.EscapeString(att.AIBOMHash) + `</code></div>
      <div>Cryptographic Signature Block: <code>` + html.EscapeString(att.SignatureBlock) + `</code></div>
    </div>
  </section>
`)
	}

	sb.WriteString(`
</div>
</body>
</html>
`)

	return sb.String()
}

func scoreColor(score float64) string {
	if score >= 80 {
		return "var(--success)"
	}
	if score >= 50 {
		return "var(--warning)"
	}
	return "var(--danger)"
}
