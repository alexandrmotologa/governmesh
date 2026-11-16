# EU AI Act Annex IV Technical Documentation Dossier

> Statutory Technical Documentation for High-Risk AI Systems pursuant to Article 11 and Annex IV of Regulation (EU) 2024/1689

**AI System / Provider:** Default Organization  
**Date of Compilation:** Tue, 17 Nov 2026 00:40:03 EET  
**Codebase Target:** `B:\workgit\governmesh\examples\enterprise_ai_service`  
**Dossier ID:** `DOC-ANNEX-IV-94710ff08c5a`  

## 1. General Description of the AI System (Annex IV.1)

### 1.1 Intended Purpose and Practical Application
This technical documentation pertains to the software components identified within `Default Organization`. The static AST inspection identified 7 integrated AI model(s) and framework components.

### 1.2 System Risk Classification
| Category | System Count | Statutory Basis |
| :--- | :--- | :--- |
| Unacceptable Risk (Article 5) | 0 | Prohibited practices (requires immediate decommissioning) |
| High-Risk (Annex III) | 3 | Subject to mandatory conformity assessment and this dossier |
| Specific Transparency (Article 50) | 1 | Subject to user disclosure notices |
| Minimal Risk (Article 69) | 3 | Standard software quality baseline |

## 2. Methods and Steps for AI Development (Annex IV.2)

### 2.1 Model Architectures and Foundation Model Supply Chain

| Component ID | Model Family | Provider / Origin | Framework | Code Location |
| :--- | :--- | :--- | :--- | :--- |
| `dock-Dockerfile-1` | **vllm/vllm-openai:latest** | Self-Hosted / Container Runtime | `vllm-container` | `Dockerfile:1` |
| `env-model-mistralai/Mistral-7B-Instruct-v0.3` | **mistralai/Mistral-7B-Instruct-v0.3** | Environment Injected Provider | `env-config` | `Dockerfile:1` |
| `prompt-hr_eval.prompt` | **Prompt Specification** | Prompt Engineering Asset | `prompt-template` | `prompts/hr_eval.prompt:1` |
| `py-app.py-1` | **gpt-4o** | OpenAI | `langchain` | `app.py:2` |
| `py-app.py-2` | **gpt-4o-mini** | OpenAI | `langchain` | `app.py:2` |
| `py-app.py-3` | **claude-3-5-sonnet** | Anthropic | `langchain` | `app.py:2` |
| `ts-edge_service.ts-1` | **gpt-4o-mini** | OpenAI | `vercel-ai/openai` | `edge_service.ts:1` |

## 3. Human Oversight & Operational Monitoring Architecture (Annex IV.3)

Pursuant to Article 14 of the EU AI Act, High-Risk AI systems must be designed to enable natural persons to oversee their operation.

- **Component `py-app.py-1` (gpt-4o):** AI used for creditworthiness evaluation, loan eligibility, or emergency service dispatch is classified High-Risk Annex III Section 5. Triggered by context identifier: 'credit_score'.
  - *Mandatory Human Review Gate:* Active CI/CD verification flag.
  - *Governing Rule:* EU AI Act Annex III Section 5 (Access to Essential Private and Public Services)
- **Component `py-app.py-2` (gpt-4o-mini):** AI used for creditworthiness evaluation, loan eligibility, or emergency service dispatch is classified High-Risk Annex III Section 5. Triggered by context identifier: 'credit_score'.
  - *Mandatory Human Review Gate:* Active CI/CD verification flag.
  - *Governing Rule:* EU AI Act Annex III Section 5 (Access to Essential Private and Public Services)
- **Component `py-app.py-3` (claude-3-5-sonnet):** AI used for creditworthiness evaluation, loan eligibility, or emergency service dispatch is classified High-Risk Annex III Section 5. Triggered by context identifier: 'credit_score'.
  - *Mandatory Human Review Gate:* Active CI/CD verification flag.
  - *Governing Rule:* EU AI Act Annex III Section 5 (Access to Essential Private and Public Services)

## 4. Data Governance and Egress Boundaries (Annex IV.4)

Pursuant to Article 10, training, validation, and testing datasets and prompt data flows must adhere to strict data governance practices.

### Detected Sensitive Variable Flows

| Variable / Snippet | PII Category | Destination Provider | Source Code Line |
| :--- | :--- | :--- | :--- |
| `{candidate_ssn}` | `pii_national_id` | Prompt Engineering Asset | `prompts/hr_eval.prompt:1` |
| `{patient_medical_history}` | `phi_medical` | Prompt Engineering Asset | `prompts/hr_eval.prompt:1` |
| `{candidate_ssn}` | `pii_national_id` | OpenAI | `app.py:1` |
| `{salary_expectation}` | `financial` | OpenAI | `app.py:1` |
| `{credit_score}` | `financial` | OpenAI | `app.py:1` |
| `{bank_account}` | `financial` | OpenAI | `app.py:1` |
| `{candidate_ssn}` | `pii_national_id` | OpenAI | `app.py:1` |
| `{salary_expectation}` | `financial` | OpenAI | `app.py:1` |
| `{credit_score}` | `financial` | OpenAI | `app.py:1` |
| `{bank_account}` | `financial` | OpenAI | `app.py:1` |
| `{candidate_ssn}` | `pii_national_id` | Anthropic | `app.py:1` |
| `{salary_expectation}` | `financial` | Anthropic | `app.py:1` |
| `{credit_score}` | `financial` | Anthropic | `app.py:1` |
| `{bank_account}` | `financial` | Anthropic | `app.py:1` |
| `{user_email}` | `pii_personal` | OpenAI | `edge_service.ts:1` |

## 5. Technical Standards and Conformity (Annex IV.5)

The system has been evaluated against harmonized standards and international governance frameworks:

- **ISO/IEC 42001:2023 Alignment:** Controls A.5, A.6, A.7, A.8, A.9 evaluated via automated AST inspection.
- **NIST AI RMF 1.0 Alignment:** Functions GOVERN, MAP, MEASURE, MANAGE verified.
- **Composite Readiness Score:** **0.0%**

## 6. Cryptographic Attestation and Signatures (Annex IV.6)

This technical dossier is cryptographically anchored to the repository source state via a SHA-256 Merkle hash chain:

- **Merkle Root Checksum:** `8ee9087c15e409a6881fbddeeba00f80da2a7d8fa0e70626d0b3624606829ef7`
- **AIBOM Checksum:** `8268eead290493086bbf8c988273decff772e669ea05138d3982c98d134a22c3`
- **Attestation Signer:** governmesh-scanner (Enterprise AI Security Office)
- **Attestation Timestamp:** 2026-11-16T22:40:03Z
- **Master Signature Block:** `659b9cf0c72d433489088235223f6642bdf06312707bef6850be5e341fb6a9fb`

---
*Dossier compiled by governmesh automated compliance engine.*  
