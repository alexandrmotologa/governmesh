import React from 'react';
import { ShieldCheck, ShieldAlert, Play, Download, RefreshCw, FileText, Code } from 'lucide-react';
import { ScanResult } from '../types';

interface HeaderProps {
  data: ScanResult | null;
  loading: boolean;
  onTriggerScan: () => void;
}

export const Header: React.FC<HeaderProps> = ({ data, loading, onTriggerScan }) => {
  return (
    <header className="border-b border-slate-800 bg-[#0d1424] px-6 py-4">
      <div className="max-w-7xl mx-auto flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div className="flex items-center gap-3">
          <div className="p-2.5 rounded-xl bg-blue-600/10 border border-blue-500/20 text-blue-400">
            <ShieldCheck className="w-6 h-6" />
          </div>
          <div>
            <div className="flex items-center gap-2">
              <h1 className="text-xl font-bold text-white tracking-tight">governmesh</h1>
              <span className="px-2 py-0.5 text-xs font-semibold rounded-full bg-blue-500/10 text-blue-400 border border-blue-500/20">
                Studio v1.0.0
              </span>
              {data && (
                <span
                  className={`px-2.5 py-0.5 text-xs font-bold rounded-full border ${
                    data.passed_quality_gate
                      ? 'bg-emerald-500/10 text-emerald-400 border-emerald-500/30'
                      : 'bg-rose-500/10 text-rose-400 border-rose-500/30'
                  }`}
                >
                  {data.passed_quality_gate ? 'CI/CD GATE: PASSED' : 'CI/CD GATE: BLOCKED'}
                </span>
              )}
            </div>
            <p className="text-xs text-slate-400 mt-0.5">
              Target: <span className="text-slate-300 font-mono">{data ? data.scanned_path : 'Repository Root'}</span>
            </p>
          </div>
        </div>

        <div className="flex items-center gap-2.5">
          <button
            onClick={onTriggerScan}
            disabled={loading}
            className="flex items-center gap-2 px-4 py-2 text-sm font-medium rounded-lg bg-blue-600 hover:bg-blue-500 text-white transition disabled:opacity-50 shadow-sm"
          >
            {loading ? <RefreshCw className="w-4 h-4 animate-spin" /> : <Play className="w-4 h-4" />}
            {loading ? 'Auditing AI Assets...' : 'Run Audit Scan'}
          </button>

          {data && (
            <div className="flex items-center gap-1.5">
              <a
                href="/api/report/html"
                target="_blank"
                rel="noreferrer"
                className="flex items-center gap-1.5 px-3 py-2 text-xs font-medium rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-200 border border-slate-700 transition"
                title="View printable executive compliance report"
              >
                <FileText className="w-3.5 h-3.5 text-slate-400" />
                Audit Report
              </a>
              <a
                href="/api/aibom"
                target="_blank"
                rel="noreferrer"
                download="aibom.json"
                className="flex items-center gap-1.5 px-3 py-2 text-xs font-medium rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-200 border border-slate-700 transition"
                title="Download CycloneDX 1.6 AIBOM JSON"
              >
                <Code className="w-3.5 h-3.5 text-slate-400" />
                CycloneDX
              </a>
            </div>
          )}
        </div>
      </div>
    </header>
  );
};
