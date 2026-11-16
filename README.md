<p align="center">
  <img src="docs/images/logo.png?raw=true" alt="governmesh Logo" width="130" style="border-radius: 24px;" />
</p>

<h1 align="center">governmesh</h1>

<p align="center">
  <strong>Automated Technical Auditing, AIBOM Generator & EU AI Act Compliance Engine in Go with Embedded React 19 Governance Studio</strong>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Language-Go%201.24-00ADD8?style=flat-square&logo=go" alt="Go" />
  <img src="https://img.shields.io/badge/Frontend-React%2019-61DAFB?style=flat-square&logo=react" alt="React 19" />
  <img src="https://img.shields.io/badge/AIBOM-CycloneDX%201.6-00f5ff?style=flat-square" alt="CycloneDX 1.6" />
  <img src="https://img.shields.io/badge/AIBOM-SPDX%203.0-4B7BE5?style=flat-square" alt="SPDX 3.0" />
  <img src="https://img.shields.io/badge/Standards-EU%20AI%20Act%20%7C%20ISO%2042001%20%7C%20NIST%20AI%20RMF-10b981?style=flat-square" alt="Standards" />
  <img src="https://img.shields.io/badge/SARIF-v2.1.0-8b5cf6?style=flat-square" alt="SARIF v2.1.0" />
  <img src="https://img.shields.io/badge/License-MIT-blue?style=flat-square" alt="License" />
</p>

<p align="center">
  <img src="docs/images/governmesh_demo.gif?raw=true" alt="governmesh Live Governance Studio & CLI Demo" width="100%" />
</p>

---

## Overview

governmesh is a single-binary static analysis engine (< 20MB RAM) with an embedded React 19 studio. It inspects Python, TypeScript, and Go codebases to identify AI model calls, prompt templates, weights artifacts, and data egress channels.

The engine classifies each discovered AI component against the 4 risk tiers of the EU AI Act (Article 5 Prohibited, Annex III High-Risk, Article 50 Specific Transparency, Article 69 Minimal Risk), generates authoritative `aibom.json` files conforming to CycloneDX 1.6 and SPDX 3.0, audits sensitive PII/PHI variable flows, and seals the audit trail in a SHA-256 Merkle hash chain.

---

## Live Studio Walkthrough

### Desktop Executive Cockpit
![governmesh Executive Cockpit](docs/images/governmesh_dashboard.png?raw=true)

### Automated Code Remediations & Patches
![governmesh Code Remediations and Patches](docs/images/governmesh_remediation.png?raw=true)

### Mobile & Responsive Viewport (390px)
<p align="center">
  <img src="docs/images/governmesh_mobile.png?raw=true" alt="governmesh Mobile Responsive View" width="380" />
</p>

### AIBOM Supply Chain Explorer
![governmesh AIBOM Explorer](docs/images/governmesh_aibom.png?raw=true)

### EU AI Act Annex IV Technical Dossier
![governmesh Annex IV Dossier](docs/images/governmesh_annex_iv.png?raw=true)

---

## Capabilities

- **Multi-Language AST Codebase Scanner**: Analyzes Python, TypeScript, JavaScript, Go, Dockerfiles, and prompt specifications for foundation model APIs, weights checkpoints, and vector integrations.
- **Automated EU AI Act Risk Tiering**: Classifies discovered models into the 4 statutory tiers: Unacceptable Risk (Article 5 Prohibitions), High-Risk (Annex III Employment, Credit, Critical Infrastructure), Specific Transparency (Article 50 Chatbots), and Minimal Risk (Article 69).
- **CycloneDX 1.6 & SPDX 3.0 AIBOM Generation**: Exports machine learning software supply chain inventories (`aibom.json`) including model cards, providers, parameters, and licensing data.
- **Sensitive Data & PII Egress Detection**: Scans prompt templates and runtime variable interpolations for SSN, IBAN, credit cards, emails, and health records sent to third-party cloud LLM APIs.
- **EU AI Act Annex IV Technical Documentation Generator**: Assembles the statutory technical dossier required under Article 11 for High-Risk AI deployments (`governmesh docgen`).
- **Automated Code Remediation & Patch Engine**: Plans targeted compliance code fixes, privacy safeguards, and Article 50 disclosure notices, generating unified git patches (`governmesh fix`).
- **SARIF v2.1.0 Export**: Produces standardized OASIS SARIF reports for native GitHub Actions security alerts and pull request line annotations (`--sarif`).
- **Cryptographic Attestation & Merkle Chains**: Computes SHA-256 hash chains across all scanned files and generated AIBOM documents, outputting a verifiable attestation manifest.
- **CI/CD Quality Gate**: Returns deterministic non-zero exit codes to block builds when prohibited practices or unmitigated high risks are detected.
- **Embedded React 19 Governance Studio**: Single-binary local console served at `http://localhost:6363` with risk matrices, interactive inventory explorer, remediation diff viewer, and Annex IV dossier generator.

---

## Installation

### Pre-built Binary

```bash
git clone https://github.com/alexandrmotologa/governmesh.git
cd governmesh
go build -o governmesh ./cmd/governmesh
```

### Go Install

```bash
go install github.com/alexandrmotologa/governmesh/cmd/governmesh@latest
```

---

## Quick Start

### 1. Run a Compliance Audit Scan

```bash
governmesh scan \
  --config examples/governmesh.yaml \
  --aibom aibom.json \
  --attestation attestation.json \
  --report compliance_report.html \
  --sarif results.sarif \
  --fail-on unacceptable \
  ./path/to/codebase
```

### 2. Generate EU AI Act Annex IV Technical Dossier

```bash
governmesh docgen --output annex_iv_dossier.md ./path/to/codebase
```

### 3. Generate Code Remediations and Git Patch

```bash
governmesh fix --output-patch remediation.patch ./path/to/codebase
```

### 4. Generate Authoritative AIBOM

```bash
# CycloneDX 1.6 format
governmesh aibom --format cyclonedx --output aibom.json ./path/to/codebase

# SPDX 3.0 format
governmesh aibom --format spdx --output aibom_spdx.json ./path/to/codebase
```

### 5. Verify Cryptographic Integrity

Inspect an audit bundle to detect if source code or AIBOM contents were modified after signing:

```bash
governmesh verify \
  --manifest attestation.json \
  --aibom aibom.json \
  --path ./path/to/codebase
```

### 6. Compare Snapshots with Diff

Detect newly introduced models or risk tier escalations between versions:

```bash
governmesh diff --base v1_aibom.json --target v2_aibom.json
```

### 7. Launch Governance Studio

Launch the embedded React 19 studio console:

```bash
governmesh studio --port 6363 ./path/to/codebase
```

Open `http://localhost:6363` in your browser.

---

## Command Reference

| Command | Purpose | Key Flags |
| :--- | :--- | :--- |
| `scan` | Complete codebase audit, tiering, and gate evaluation | `--config`, `--aibom`, `--attestation`, `--report`, `--sarif`, `--docgen`, `--fail-on`, `--json` |
| `docgen` | Generate statutory EU AI Act Annex IV technical dossier | `--output`, `--config` |
| `fix` | Plan compliance code remediations and produce git patch | `--output-patch`, `--dry-run`, `--config` |
| `aibom` | Synthesize authoritative AIBOM JSON | `--format` (cyclonedx, spdx), `--output` |
| `verify` | Verify cryptographic SHA-256 Merkle chain | `--manifest`, `--aibom`, `--path` |
| `report` | Generate executive audit document | `--format` (html, md), `--output`, `--path` |
| `diff` | Analyze risk delta between two AIBOM files | `--base`, `--target` |
| `studio` | Launch embedded React 19 web dashboard | `--port` (default 6363), `--config` |
| `version` | Display version and compliance engine details | |

---

## Supported Standards

- EU AI Act risk tiers and prohibited practices
- NIST AI Risk Management Framework functions
- ISO 42001 management controls

---

## Architecture & Documentation

- [System Architecture](docs/ARCHITECTURE.md)
- [EU AI Act Mapping Taxonomy](docs/EU_AI_ACT_MAPPING.md)
- [Standards Alignment (ISO 42001 & NIST AI RMF)](docs/STANDARDS.md)
- [CI/CD Quality Gate Integration](docs/CI_CD_INTEGRATION.md)
- [AIBOM Specification (CycloneDX 1.6 & SPDX 3.0)](docs/AIBOM_SPECIFICATION.md)

---

## License

MIT License. See [LICENSE](LICENSE) for details.
