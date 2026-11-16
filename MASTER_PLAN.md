# governmesh: Architecture and Engineering Master Plan
> Enterprise AI Bill of Materials (AIBOM) & Automated EU AI Act / ISO 42001 Compliance Scanner
> Target: Single-Binary Go Engine (< 20MB RAM) with Embedded React 19 Governance Studio & CycloneDX AI Exporter

---

## 1. Problem Statement and Motivation

### 1.1 The Regulatory Pressure on Enterprise AI Deployments
With the full enforcement of the European Union Artificial Intelligence Act (EU AI Act), the NIST AI Risk Management Framework (AI RMF), and the ISO/IEC 42001 standard, enterprise AI governance has transitioned from voluntary best practices to strict legal mandates with severe non-compliance penalties (up to €35 million or 7% of global annual turnover):
1. **The Proliferation of Shadow AI**: Developers and product teams rapidly integrate foundation model APIs (OpenAI, Anthropic, Bedrock, Vertex), open-source weights (Hugging Face), and vector databases without centralized tracking. Most enterprises cannot produce a comprehensive inventory of where AI is embedded across their codebase.
2. **Mandatory AI Bill of Materials (AIBOM)**: Regulatory compliance requires organizations to document the complete supply chain of every AI application: model architectures, pre-training cutoffs, dataset provenance, fine-tuning pipelines, system prompt invariants, and third-party data egress points.
3. **Subjective Risk Classification**: The EU AI Act enforces asymmetric legal obligations based on risk tiers (Unacceptable, High-Risk, Specific Transparency, and Minimal Risk). Misclassifying a CV screening or credit assessment tool as "Minimal Risk" triggers severe statutory liability.
4. **Manual, Static Spreadsheet Audits**: Compliance teams currently attempt to audit AI assets using static Excel questionnaires that are outdated before completion, failing to capture continuous CI/CD deployments and runtime model changes.

### 1.2 The governmesh Solution
`governmesh` is an automated AI governance and compliance platform written in Go. Operating as a single-binary CLI and CI/CD quality gate:
- **Multi-Language AST Codebase Scanner**: Scans Git repositories (Python, TypeScript, Go, Java, Docker, Kubernetes manifests) to discover all LLM integrations, model weights, API clients, vector stores, and prompt templates.
- **Automated EU AI Act Risk Tiering**: Classifies each discovered AI workflow against the EU AI Act 4-tier risk taxonomy and NIST AI RMF categories using deterministic heuristics and keyword intent mapping.
- **CycloneDX 1.6 & SPDX AI AIBOM Generation**: Employs open standards to export an authoritative `aibom.json` detailing model versions, parameter sizes, license terms, dataset hashes, and API endpoints.
- **PII and Sensitive Data Egress Detection**: Audits prompt assembly templates and RAG retrieval pipelines for potential exfiltration of personal identifiers (PII), protected health information (PHI), or financial credentials.
- **Cryptographic Compliance Attestation**: Signs the audit trail with SHA-256 hash chains, exporting legal-grade compliance reports (PDF and Markdown) ready for internal counsel and external regulatory inspection.

---

## 2. High-Level System Architecture

```
       +-------------------------------------------------------------+
       |             Enterprise Software Repositories & CI/CD        |
       |  (Python, TypeScript, Go, Dockerfiles, K8s, Prompt Yamls)   |
       +-------------------------------------------------------------+
                                      |
                         Static AST & Manifest Scanning
                                      v
       +-------------------------------------------------------------+
       |                     governmesh Core Engine                  |
       |                   (High-Performance Go CLI)                |
       |                                                             |
       |  +---------------------+        +------------------------+  |
       |  | Multi-Language AST  |------->| AI Asset Inventory     |  |
       |  | Codebase Scanner    |        | Extractor              |  |
       |  +---------------------+        +------------------------+  |
       |                                              |              |
       |                                              v              |
       |  +-------------------------------------------------------+  |
       |  |           EU AI Act & NIST Risk Classifier            |  |
       |  |  - Unacceptable Risk (Prohibited Social Scoring)     |  |
       |  |  - High-Risk (Hiring, Credit, Critical Infra)         |  |
       |  |  - Specific Transparency (Chatbots, Synthetic Media)  |  |
       |  |  - Minimal Risk (Spam Filters, Search Ranking)        |  |
       |  +-------------------------------------------------------+  |
       |                               |                             |
       |                               v                             |
       |  +-------------------------------------------------------+  |
       |  |              Data Lineage & PII Egress Auditor        |  |
       |  |  - Scans Prompt Context for PII/PHI Injection Risk    |  |
       |  |  - Analyzes Vector Store Indexing and Source Docs     |  |
       |  +-------------------------------------------------------+  |
       |                               |                             |
       |                               v                             |
       |  +-------------------------------------------------------+  |
       |  |           CycloneDX 1.6 / SPDX AIBOM Synthesizer      |  |
       |  |  - Authoritative `aibom.json` Software Supply Chain   |  |
       |  +-------------------------------------------------------+  |
       |                               |                             |
       +-------------------------------|-----------------------------+
                                       |
                                       v
       +-------------------------------------------------------------+
       |                  Output & Attestation Surfaces              |
       |  1. Embedded React 19 Studio (Local HTTP :6363)             |
       |  2. Legal-Grade Compliance Report (Markdown & PDF)          |
       |  3. CI/CD Gate Exit Code (Fail on Unacceptable/High Risk)   |
       +-------------------------------------------------------------+
```

---

## 3. Regulatory Classification Framework (EU AI Act & NIST AI RMF)

```
+-----------------------------------+-----------------------------------------------------------------+
| Risk Category (EU AI Act)         | Enterprise Criteria & Action Mandate                            |
+-----------------------------------+-----------------------------------------------------------------+
| 1. Unacceptable Risk              | Prohibited: Cognitive manipulation, social scoring, biometric   |
|                                   | categorization, real-time remote biometric identification.      |
|                                   | Action: CI/CD HARD BLOCK (Exit Code 1).                         |
| 2. High-Risk                      | Permitted with strict conformity: CV screening, employment      |
|                                   | decisions, credit underwriting, law enforcement, healthcare.    |
|                                   | Action: Requires Risk Management System & Human Oversight Plan. |
| 3. Specific Transparency Risk     | Permitted with user disclosure: Customer-facing chatbots,       |
|                                   | synthetic audio/video generation, automated text generators.    |
|                                   | Action: Enforces mandatory user-facing AI disclosure badge.     |
| 4. Minimal / Low Risk             | Permitted freely: Internal code completion, search indexing,    |
|                                   | spam classification, automated spell-checking.                  |
|                                   | Action: Documented in AIBOM with standard security tracking.    |
+-----------------------------------+-----------------------------------------------------------------+
```

---

## 4. Repository Structure

```
governmesh/
├── cmd/
│   └── governmesh/
│       └── main.go                 # CLI entry point, subcommands: scan, aibom, verify, studio
├── internal/
│   ├── config/
│   │   ├── config.go               # Configuration parser for compliance rules and thresholds
│   │   └── schema.go               # Custom enterprise organizational units and risk exceptions
│   ├── scanner/
│   │   ├── ast_python.go           # Python AST parser (langchain, openai, anthropic imports)
│   │   ├── ast_typescript.go       # TypeScript/JavaScript AST parser (Vercel AI SDK, LangChainJS)
│   │   ├── ast_go.go               # Go AST parser for native Go AI SDKs
│   │   ├── manifest.go             # Dockerfile, Helm chart, and Kubernetes model volume detector
│   │   └── prompt_scanner.go       # Discovers hardcoded system prompts and prompt templates
│   ├── classifier/
│   │   ├── eu_ai_act.go            # Rule engine mapping use-cases to EU AI Act risk tiers
│   │   ├── nist_rmf.go             # NIST AI Risk Management Framework compliance evaluator
│   │   └── iso_42001.go            # ISO/IEC 42001 controls checklist mapper
│   ├── pii/
│   │   ├── egress_scanner.go       # Scans prompt templates for PII/PHI variable bindings
│   │   └── patterns.go             # Pre-compiled regex patterns for SSN, IBAN, emails, health terms
│   ├── aibom/
│   │   ├── cyclonedx.go            # CycloneDX 1.6 AI Bill of Materials JSON generator
│   │   ├── spdx.go                 # SPDX 3.0 AI extension format serializer
│   │   └── schema.go               # Internal standard AIBOM data representation
│   ├── attestation/
│   │   ├── signer.go               # Cryptographic SHA-256 signature and timestamp generator
│   │   └── pdf_report.go           # Executive PDF audit summary generator
│   └── server/
│       ├── api.go                  # REST API for compliance status and repository scans
│       ├── embed.go                # Embedded React 19 studio static assets
│       └── sse.go                  # Real-time scan progress event stream
├── web/
│   ├── src/
│   │   ├── App.tsx                 # React 19 root application
│   │   ├── components/
│   │   │   ├── RiskMatrixChart.tsx # 4-tier EU AI Act distribution donut and risk heatmap
│   │   │   ├── AibomExplorer.tsx   # Interactive tree view of models, datasets, and endpoints
│   │   │   ├── PiiEgressTable.tsx  # Table of discovered prompt PII risks and file locations
│   │   │   ├── ComplianceScore.tsx # Percentage compliance readiness gauge (ISO 42001 / NIST)
│   │   │   └── ExportButton.tsx    # One-click CycloneDX JSON & Signed PDF export
│   │   └── pages/
│   │       └── Dashboard.tsx       # Main executive compliance cockpit
│   ├── package.json
│   ├── tsconfig.json
│   └── vite.config.ts
├── examples/
│   ├── governmesh.yaml             # Example compliance configuration
│   └── sample_repo/                # Mock repository with diverse AI integrations for scanning
│       ├── app.py
│       ├── prompt_templates.yaml
│       └── Dockerfile
├── Makefile
├── go.mod
├── go.sum
├── LICENSE
└── README.md
```

---

## 5. Implementation Roadmap

### Phase 1: Go Engine Setup and Multi-Language AST Scanner
- Initialize Go module with minimal external dependencies.
- Build AST parsers for Python (`internal/scanner/ast_python.go`) and TypeScript/JavaScript detecting imports from `openai`, `@anthropic-ai/sdk`, `ai` (Vercel), and `@langchain/core`.
- Scan environment files and Kubernetes manifests for model endpoint secrets and configurations.

### Phase 2: Regulatory Risk Classification Engine
- Build `internal/classifier/eu_ai_act.go` implementing the 4-tier risk decision tree.
- Identify sensitive domain indicators (e.g. CV analysis, credit evaluation, facial recognition) via AST variable names and function signatures.
- Implement NIST AI RMF governance controls evaluator.

### Phase 3: PII and Data Egress Auditor
- Implement `internal/pii` scanning prompt template interpolations (e.g. `{customer_ssn}`, `{medical_history}`).
- Flag unmasked sensitive attributes sent to external public LLM API endpoints.

### Phase 4: CycloneDX 1.6 AIBOM Generator
- Implement `internal/aibom/cyclonedx.go` generating standard JSON AIBOM containing components, model specifications, training dataset citations, and licensing metadata.
- Support SPDX 3.0 serialization.

### Phase 5: Cryptographic Attestation and CI/CD Quality Gates
- Implement SHA-256 hashing across all scanned artifacts, generating an immutable compliance manifest.
- Support CLI exit codes (`0` for passing, `1` if Unacceptable Risk or Unmitigated High Risk is detected in CI).
- Generate legal-grade executive Markdown and PDF compliance summaries.

### Phase 6: Embedded React 19 Governance Studio
- Develop compliance cockpit in `web/` featuring risk distribution charts, AIBOM tree viewer, and PII risk tracker.
- Embed frontend static assets into the Go binary via `//go:embed` and serve locally on `:6363`.

### Phase 7: Verification, Benchmarking, and Documentation
- Run test scans on `examples/sample_repo` verifying accurate detection and classification of all AI assets in $< 500\text{ms}$.
- Verify binary memory usage stays strictly under 20MB RAM.
- Author humanized README with setup guides, EU AI Act compliance checklists, and CI/CD pipeline integration examples.
