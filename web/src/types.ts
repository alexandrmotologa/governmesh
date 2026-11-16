export type RiskLevel = 'unacceptable' | 'high' | 'transparency' | 'minimal';

export interface PIIFinding {
  type: string;
  snippet: string;
  line: number;
  severity: string;
  variable?: string;
  description: string;
}

export interface RiskAssessment {
  tier: RiskLevel;
  applicable_article: string;
  justification: string;
  prohibited_reason?: string;
  requires_human_gate: boolean;
  nist_categories: string[];
  iso_42001_controls: string[];
}

export interface DiscoveredAsset {
  id: string;
  name: string;
  model_family: string;
  provider: string;
  source_file_path: string;
  line_number: number;
  language: string;
  framework: string;
  context_snippet: string;
  parameters?: Record<string, unknown>;
  system_prompt?: string;
  prompt_vars?: string[];
  egress_target?: string;
  pii_findings?: PIIFinding[];
  intent_keywords?: string[];
  risk: RiskAssessment;
  file_hash: string;
}

export interface PIIExposure {
  asset_id: string;
  file_path: string;
  line: number;
  pii_type: string;
  target_cloud: string;
  snippet: string;
  is_prohibited: boolean;
}

export interface ScanResult {
  project_name: string;
  scanned_path: string;
  scan_timestamp: string;
  duration_ms: number;
  total_files: number;
  scanned_ai_files: number;
  assets: DiscoveredAsset[];
  risk_summary: Record<RiskLevel, number>;
  pii_exposures: PIIExposure[];
  compliance_score: number;
  hash_chain: string;
  passed_quality_gate: boolean;
  gate_failures?: string[];
}

export interface AttestationManifest {
  version: string;
  project_name: string;
  timestamp: string;
  signer: string;
  organization: string;
  algorithm: string;
  root_hash: string;
  aibom_hash: string;
  file_tree_hashes: Record<string, string>;
  risk_signature: string;
  signature_block: string;
}
