import React from 'react';
import { ShieldAlert, AlertCircle, CheckCircle, ExternalLink } from 'lucide-react';
import { PIIExposure } from '../types';

interface PiiEgressTableProps {
  exposures: PIIExposure[];
}

export const PiiEgressTable: React.FC<PiiEgressTableProps> = ({ exposures }) => {
  return (
    <div className="bg-[#0f172a] rounded-xl border border-slate-800 overflow-hidden">
      <div className="p-4 border-b border-slate-800 flex justify-between items-center">
        <div>
          <h3 className="text-sm font-semibold text-white flex items-center gap-2">
            <ShieldAlert className="w-4 h-4 text-amber-400" />
            Sensitive Data & PII Egress Audit
          </h3>
          <p className="text-xs text-slate-400 mt-0.5">
            Personal identifiers or sensitive demographic data mapped into prompt templates and egressing to external cloud LLMs.
          </p>
        </div>
        <span className="px-2.5 py-1 rounded-full text-xs font-bold bg-slate-800 text-slate-300">
          {exposures.length} Total Findings
        </span>
      </div>

      <div className="overflow-x-auto">
        <table className="w-full text-left border-collapse">
          <thead>
            <tr className="bg-slate-900/60 text-slate-400 text-xs font-semibold uppercase tracking-wider border-b border-slate-800">
              <th className="py-3 px-4">Severity & Status</th>
              <th className="py-3 px-4">PII Category</th>
              <th className="py-3 px-4">Target Cloud Provider</th>
              <th className="py-3 px-4">Interpolated Variable / Snippet</th>
              <th className="py-3 px-4">Source Code Location</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-800 text-xs">
            {exposures.length === 0 ? (
              <tr>
                <td colSpan={5} className="py-8 text-center text-slate-500">
                  <div className="flex flex-col items-center justify-center gap-1">
                    <CheckCircle className="w-6 h-6 text-emerald-500 mb-1" />
                    <span className="text-slate-300 font-medium">Clean Data Perimeter</span>
                    <span>Zero sensitive PII or PHI variables detected egressing to external LLM providers.</span>
                  </div>
                </td>
              </tr>
            ) : (
              exposures.map((exp, idx) => (
                <tr key={idx} className="hover:bg-slate-800/30 transition">
                  <td className="py-3 px-4">
                    {exp.is_prohibited ? (
                      <span className="px-2 py-0.5 rounded text-[11px] font-bold bg-rose-500/20 text-rose-400 border border-rose-500/30">
                        CRITICAL VIOLATION
                      </span>
                    ) : (
                      <span className="px-2 py-0.5 rounded text-[11px] font-bold bg-amber-500/20 text-amber-400 border border-amber-500/30">
                        REVIEW REQUIRED
                      </span>
                    )}
                  </td>
                  <td className="py-3 px-4 font-mono text-slate-200">{exp.pii_type}</td>
                  <td className="py-3 px-4 text-slate-300">{exp.target_cloud}</td>
                  <td className="py-3 px-4">
                    <code className="text-amber-300 bg-amber-950/40 px-2 py-0.5 rounded border border-amber-900/50">
                      {exp.snippet}
                    </code>
                  </td>
                  <td className="py-3 px-4 text-slate-400 font-mono text-[11px]">
                    {exp.file_path}:{exp.line}
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
};
