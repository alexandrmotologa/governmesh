# Standards Alignment: ISO/IEC 42001 and NIST AI RMF

governmesh maps technical AST discoveries, configuration files, and data flow patterns directly to recognized international standards for Artificial Intelligence Management Systems.

## ISO/IEC 42001:2023 (Annex A Controls)

| Clause | Control Name | governmesh Technical Implementation |
| :--- | :--- | :--- |
| **A.5.1** | Policies for AI Development | Enforces enterprise risk policy in configuration (`governmesh.yaml`) |
| **A.6.1** | AI Risk Assessment Process | Deterministic evaluation against EU AI Act Article 5 and Annex III taxonomies |
| **A.7.2** | Data for AI Systems & Privacy | AST prompt analysis for unmasked PII/PHI variable bindings egressing to cloud APIs |
| **A.8.2** | AI System Inventory & Traceability | Automated generation of CycloneDX 1.6 AIBOM with SHA-256 Merkle chain attestation |
| **A.8.4** | Human Oversight for High-Risk Systems | Quality gate flags Annex III components requiring mandatory human review checkpoints |
| **A.9.1** | Transparency and Disclosure | Tracks Article 50 disclosure obligations for customer-facing chatbots and synthetic media |

## NIST AI Risk Management Framework (AI RMF 1.0)

| Function | Subcategory | governmesh Technical Implementation |
| :--- | :--- | :--- |
| **GOVERN** | **GOVERN-1.1** | Statutory AI obligations mapped to technical components in code |
| **GOVERN** | **GOVERN-2.1** | Authoritative model inventory tracking versions, providers, and source locations |
| **MAP** | **MAP-1.1** | Deployment context and third-party API dependencies extracted from AST |
| **MAP** | **MAP-1.5** | High-risk domain indicators detected across function names, classes, and variable bindings |
| **MEASURE** | **MEASURE-1.1** | Sensitive data perimeter scan evaluating PII/PHI interpolation in prompt templates |
| **MANAGE** | **MANAGE-1.3** | Third-party foundation model dependencies tracked with external references in SBOM |
| **MANAGE** | **MANAGE-2.4** | Automated CI/CD gate blocking unauthorized or unmitigated deployments |
