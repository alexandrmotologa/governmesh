# CI/CD Quality Gate Integration

governmesh functions as an automated compliance quality gate within CI/CD pipelines. It scans every commit and pull request, produces an authoritative AIBOM, checks cryptographic integrity, and enforces risk thresholds.

## Exit Codes

| Code | Meaning | Condition |
| :--- | :--- | :--- |
| `0` | Success | All compliance thresholds met; zero blocking violations |
| `1` | Failure | Prohibited/Unacceptable risk detected, high-risk threshold breached, or critical PII egress detected |

## GitHub Actions

Create `.github/workflows/governmesh.yml`:

```yaml
name: AI Governance & Compliance Gate

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  compliance-audit:
    runs-on: ubuntu-latest

    steps:
      - name: Checkout Codebase
        uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.24'

      - name: Install governmesh
        run: |
          go install github.com/alexandrmotologa/governmesh/cmd/governmesh@latest

      - name: Execute governmesh Audit Gate
        run: |
          governmesh scan \
            --config .governmesh.yaml \
            --aibom aibom.json \
            --attestation attestation.json \
            --report compliance_report.html \
            --fail-on unacceptable \
            .

      - name: Archive Compliance Artifacts
        if: always()
        uses: actions/upload-artifact@v4
        with:
          name: ai-compliance-artifacts
          path: |
            aibom.json
            attestation.json
            compliance_report.html
```

## GitLab CI

Create `.gitlab-ci.yml`:

```yaml
stages:
  - compliance

governmesh-audit:
  stage: compliance
  image: golang:1.24
  script:
    - go install github.com/alexandrmotologa/governmesh/cmd/governmesh@latest
    - governmesh scan --aibom aibom.json --attestation attestation.json --report compliance_report.html --fail-on unacceptable .
  artifacts:
    when: always
    paths:
      - aibom.json
      - attestation.json
      - compliance_report.html
```

## Regression Detection with `diff`

Detect newly introduced models or risk tier escalations between release branches:

```bash
# Scan previous release
git checkout tags/v1.0.0
governmesh aibom --output base_aibom.json .

# Scan current branch
git checkout main
governmesh aibom --output target_aibom.json .

# Run diff inspection
governmesh diff --base base_aibom.json --target target_aibom.json
```
