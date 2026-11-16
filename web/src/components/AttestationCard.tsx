import React, { useState } from 'react';
import { KeyRound, CheckCircle2, ShieldCheck, Copy, Check } from 'lucide-react';
import { ScanResult } from '../types';

interface AttestationCardProps {
  data: ScanResult;
}

export const AttestationCard: React.FC<AttestationCardProps> = ({ data }) => {
  const [copied, setCopied] = useState(false);

  const copyHash = () => {
    navigator.clipboard.writeText(data.hash_chain);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div className="bg-[#0f172a] rounded-xl border border-slate-800 p-5">
      <div className="flex items-center justify-between mb-4">
        <div className="flex items-center gap-2">
          <KeyRound className="w-5 h-5 text-purple-400" />
          <h3 className="text-sm font-semibold text-white">Cryptographic Audit Attestation</h3>
        </div>
        <span className="flex items-center gap-1 text-xs text-emerald-400 font-semibold bg-emerald-500/10 px-2 py-0.5 rounded border border-emerald-500/20">
          <CheckCircle2 className="w-3.5 h-3.5" />
          SHA-256 Merkle Chain Sealed
        </span>
      </div>

      <p className="text-xs text-slate-400 leading-relaxed mb-4">
        Every scanned AST node, model call signature, and configuration parameter is bound into an immutable
        SHA-256 Merkle hash chain. Regulatory inspectors and internal counsel can independently verify audit
        authenticity and detect subsequent source tampering.
      </p>

      <div className="space-y-3 font-mono text-xs">
        <div>
          <div className="text-slate-500 text-[11px] mb-1">Merkle Root Hash Chain:</div>
          <div className="flex items-center justify-between p-2.5 rounded bg-slate-950 border border-slate-800 text-slate-300">
            <span className="truncate pr-2">{data.hash_chain || 'Generating...'}</span>
            <button
              onClick={copyHash}
              className="text-slate-400 hover:text-white transition flex items-center gap-1 text-[11px] flex-shrink-0"
            >
              {copied ? <Check className="w-3.5 h-3.5 text-emerald-400" /> : <Copy className="w-3.5 h-3.5" />}
              {copied ? 'Copied' : 'Copy'}
            </button>
          </div>
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 text-slate-400 text-[11px]">
          <div>
            <span className="text-slate-500">Audit Timestamp:</span>{' '}
            <span className="text-slate-300">{new Date(data.scan_timestamp).toUTCString()}</span>
          </div>
          <div>
            <span className="text-slate-500">Scan Duration:</span>{' '}
            <span className="text-slate-300">{data.duration_ms} ms</span>
          </div>
        </div>
      </div>
    </div>
  );
};
