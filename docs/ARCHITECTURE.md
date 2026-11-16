# governmesh Technical Architecture

governmesh is a single-binary static analysis engine and compliance platform written in Go, bundled with an embedded React 19 Governance Studio. It automates technical audits and inventorying of enterprise AI systems according to the EU AI Act, NIST AI RMF, and ISO/IEC 42001.

## System Components

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
|  |  - Unacceptable Risk (Article 5 Prohibitions)         |  |
|  |  - High-Risk (Annex III Employment, Credit, etc.)     |  |
|  |  - Specific Transparency (Article 50 Chatbots)        |  |
|  |  - Minimal Risk (Developer Tools, Search Indexing)    |  |
|  +-------------------------------------------------------+  |
|                               |                             |
|                               v                             |
|  +-------------------------------------------------------+  |
|  |              Data Lineage & PII Egress Auditor        |  |
|  |  - Scans Prompt Context for PII/PHI Injection Risk    |  |
|  |  - Detects External Cloud vs Local Inference Egress   |  |
|  +-------------------------------------------------------+  |
|                               |                             |
|                               v                             |
|  +-------------------------------------------------------+  |
|  |           CycloneDX 1.6 / SPDX AIBOM Synthesizer      |  |
|  |  - Authoritative aibom.json Software Supply Chain     |  |
|  +-------------------------------------------------------+  |
|                               |                             |
+-------------------------------|-----------------------------+
                                |
                                v
+-------------------------------------------------------------+
|                  Output & Attestation Surfaces              |
|  1. Embedded React 19 Studio (Local HTTP :6363)             |
|  2. Legal-Grade Compliance Report (Markdown & HTML)         |
|  3. CI/CD Gate Exit Code (Fail on Unacceptable/High Risk)   |
|  4. SHA-256 Merkle Chain Attestation Manifest               |
+-------------------------------------------------------------+
```

## Core Subsystems

### 1. Multi-Language AST Codebase Scanner
- **Python AST Parser**: Identifies imports and invocations of OpenAI, Anthropic, LangChain, Transformers, LiteLLM, vLLM, and Ollama. Extracts prompt templates, system instructions, and variables.
- **TypeScript and JavaScript Parser**: Detects Vercel AI SDK, LangChainJS, OpenAI, and Anthropic SDKs. Analyzes template literals and variable interpolation.
- **Go AST Parser**: Uses the Go standard library (`go/parser` and `go/ast`) to detect Go-based AI libraries such as `go-openai` and `langchaingo`.
- **Infrastructure Manifest Parser**: Scans Dockerfiles, container base images (`ollama/ollama`, `vllm/vllm-openai`), Kubernetes manifests, and `.env` files for embedded model checkpoints (`.safetensors`, `.onnx`, `.gguf`) and credentials.
- **Prompt Specification Parser**: Scans dedicated prompt template files (`.prompt`, `.yaml`, `.json`).

### 2. Regulatory Risk Classifier
Evaluates every discovered model call against:
- **EU AI Act**: Article 5 (Prohibited practices), Article 6 & Annex III (High-risk areas: recruitment, credit, critical infrastructure), Article 50 (Transparency obligations for chatbots), and Article 69 (Minimal risk).
- **NIST AI RMF 1.0**: Maps to GOVERN, MAP, MEASURE, and MANAGE functions.
- **ISO/IEC 42001:2023**: Maps controls across Annex A clauses (A.5 through A.9).

### 3. PII and Data Egress Inspector
Audits variable bindings and prompt literals for sensitive personal information:
- SSN, European National IDs, Romanian CNP
- Credit cards and IBANs
- Email addresses and phone numbers
- Protected Health Information (PHI) keywords
- Financial variables (credit scores, bank account numbers, salary)
- Special category demographics (religion, political beliefs, biometric data)

### 4. Authoritative AIBOM Synthesizer
- Generates official CycloneDX 1.6 JSON with machine-learning-model components, model cards, properties, and external references.
- Exports SPDX 3.0 AI profile elements and relationships.

### 5. Cryptographic Attestation Engine
- Computes SHA-256 Merkle root across all scanned files and generated AIBOM.
- Produces an immutable `attestation.json` manifest.
- Includes a verification engine (`governmesh verify`) that detects any file tampering or post-audit alterations.
