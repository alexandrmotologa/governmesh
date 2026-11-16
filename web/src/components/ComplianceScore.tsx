import React from 'react';
import { Award, ShieldCheck, ShieldAlert, Check, X } from 'lucide-react';
import { ScanResult } from '../types';

interface ComplianceScoreProps {
  data: ScanResult;
}

export const ComplianceScore: React.FC<ComplianceScoreProps> = ({ data }) => {
  const score = data.compliance_score;
  const scoreColor =
    score >= 80 ? 'text-emerald-400' : score >= 50 ? 'text-amber-400' : 'text-rose-400';
  const strokeColor =
    score >= 80 ? '#10b981' : score >= 50 ? '#f59e0b' : '#ef4444';

  return (
    <div className="bg-[#0f172a] p-5 rounded-xl border border-slate-800 flex flex-col justify-between">
      <div>
        <div className="flex items-center justify-between mb-4">
          <div className="flex items-center gap-2">
            <Award className="w-5 h-5 text-blue-400" />
            <h3 className="text-sm font-semibold text-white">Compliance Readiness</h3>
          </div>
          <span className="text-xs px-2 py-0.5 rounded bg-slate-800 text-slate-300 font-mono">
            ISO 42001 & NIST AI RMF
          </span>
        </div>

        <div className="flex items-center gap-6 my-2">
          {/* Circular progress visual */}
          <div className="relative w-20 h-20 flex-shrink-0 flex items-center justify-center">
            <svg className="w-20 h-20 transform -rotate-90">
              <circle
                cx="40"
                cy="40"
                r="32"
                stroke="#1e293b"
                strokeWidth="6"
                fill="transparent"
              />
              <circle
                cx="40"
                cy="40"
                r="32"
                stroke={strokeColor}
                strokeWidth="6"
                strokeDasharray={201}
                strokeDashoffset={201 - (201 * score) / 100}
                strokeLinecap="round"
                fill="transparent"
                className="transition-all duration-700 ease-out"
              />
            </svg>
            <div className="absolute text-center">
              <span className={`text-xl font-black ${scoreColor}`}>{score.toFixed(0)}</span>
              <span className="text-[10px] text-slate-400 font-semibold">%</span>
            </div>
          </div>

          <div className="space-y-1 text-xs">
            <div className="text-slate-300 font-medium">Composite AI Governance Index</div>
            <p className="text-slate-400 leading-relaxed">
              Synthesized from static AST evidence, sensitive data egress audits, and mandatory
              controls mapping.
            </p>
          </div>
        </div>
      </div>

      <div className="mt-4 pt-3 border-t border-slate-800/80 space-y-2">
        <div className="flex items-center justify-between text-xs">
          <span className="text-slate-400">CI/CD Quality Gate</span>
          <span
            className={`font-semibold ${
              data.passed_quality_gate ? 'text-emerald-400' : 'text-rose-400'
            }`}
          >
            {data.passed_quality_gate ? 'PASSED (0 Blocker Violations)' : 'FAILED (Remediation Required)'}
          </span>
        </div>
        {data.gate_failures && data.gate_failures.length > 0 && (
          <div className="text-[11px] text-rose-300/90 bg-rose-950/30 p-2 rounded border border-rose-900/50">
            {data.gate_failures[0]}
          </div>
        )}
      </div>
    </div>
  );
};
