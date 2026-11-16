# AI Bill of Materials (AIBOM) Specification

governmesh generates authoritative AI Bill of Materials documents compliant with open supply chain standards: CycloneDX 1.6 and SPDX 3.0.

## 1. CycloneDX 1.6 Format

CycloneDX 1.6 introduces native extensions for machine learning models and dataset provenance.

### Top-Level Document Structure
```json
{
  "$schema": "http://cyclonedx.org/schema/bom-1.6.schema.json",
  "bomFormat": "CycloneDX",
  "specVersion": "1.6",
  "serialNumber": "urn:uuid:852df244-f743-4e6e-9726-a6d481a451fd",
  "version": 1,
  "metadata": {
    "timestamp": "2026-11-16T22:28:05Z",
    "tools": {
      "components": [
        {
          "type": "application",
          "name": "governmesh",
          "version": "1.0.0"
        }
      ]
    },
    "component": {
      "type": "application",
      "name": "Enterprise Core AI Services",
      "version": "1.0.0"
    },
    "properties": [
      { "name": "governmesh:compliance:score", "value": "85.0" },
      { "name": "governmesh:compliance:passed_gate", "value": "true" },
      { "name": "governmesh:attestation:hash_chain", "value": "1c2e253fbcd1ede6aaad4307bbd00a85358509734c637614e86b047e4d193f89" }
    ]
  },
  "components": [
    {
      "type": "machine-learning-model",
      "name": "gpt-4o",
      "version": "1.0.0",
      "purl": "pkg:ml/openai/gpt-4o@latest",
      "modelCard": {
        "modelParameters": {
          "task": "text-and-multimodal-generation",
          "architectureFamily": "transformer"
        },
        "modelConsiderations": {
          "users": ["enterprise-application"],
          "useCases": ["EU AI Act Annex III Section 4 (Employment, Workers Management)"]
        },
        "properties": [
          { "name": "governmesh:risk:tier", "value": "high" },
          { "name": "governmesh:egress:target", "value": "external_cloud" },
          { "name": "governmesh:source:file", "value": "app.py" },
          { "name": "governmesh:source:line", "value": "10" }
        ]
      }
    }
  ]
}
```

## 2. SPDX 3.0 Format

governmesh also supports SPDX 3.0 serialization with `AIModel` elements and dependency relationships:

```json
{
  "spdxVersion": "SPDX-3.0",
  "dataLicense": "CC0-1.0",
  "spdxId": "SPDXRef-DOCUMENT",
  "name": "AIBOM-Enterprise Core AI Services",
  "elements": [
    {
      "type": "SoftwarePackage",
      "spdxId": "SPDXRef-RootPackage",
      "name": "Enterprise Core AI Services"
    },
    {
      "type": "AIModel",
      "spdxId": "SPDXRef-AIModel-1",
      "name": "gpt-4o",
      "originator": "OpenAI",
      "safetyRiskCategory": "high"
    }
  ],
  "relationships": [
    {
      "from": "SPDXRef-RootPackage",
      "to": "SPDXRef-AIModel-1",
      "relationshipType": "DEPENDS_ON"
    }
  ]
}
```
