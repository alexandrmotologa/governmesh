import React, { useState, useEffect } from 'react';
import { Wrench, Check, Copy, FileCode, ShieldAlert, Sparkles } from 'lucide-react';
import { ProposedFix, RemediationPayload } from '../types';

export const RemediationViewer: React.FC = () => {
  const [data, setData] = useState<RemediationPayload | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [copied, setCopied] = useState<boolean>(false);

  useEffect(() => {
    fetch('/api/remediation')
      .then((res) => res.json())
      .then((payload: RemediationPayload) => {
        setData(payload);
        setLoading(false);
      })
      .catch(() => setLoading(false));
  }, []);

  const copyDiff = () => {
    if (data?.diff) {
      navigator.clipboard.writeText(data.diff);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    }
  };

  if (loading) {
    return (
      <div className="bg-[#0f172a] rounded-xl border border-slate-800 p-8 text-center text-slate-400 text-xs">
        Synthesizing automated compliance remediations...
      </div>
    );
  }

  const fixes = data?.fixes || [];

  return (
    <div className="bg-[#0f172a] rounded-xl border border-slate-800 overflow-hidden space-y-4">
      <div className="p-4 border-b border-slate-800 flex justify-between items-center">
        <div>
          <h3 className="text-sm font-semibold text-white flex items-center gap-2">
            <Wrench className="w-4 h-4 text-emerald-400" />
            Automated Compliance Remediations & Patches
          </h3>
          <p className="text-xs text-slate-400 mt-0.5">
            Auto-generated patches: Article 50 transparency notices, Article 14 human oversight hooks, and sensitive variable safeguards.
          </p>
        </div>
        <div className="flex items-center gap-2">
          <span className="px-2.5 py-1 rounded-full text-xs font-bold bg-emerald-500/10 text-emerald-400 border border-emerald-500/20">
            {fixes.length} Planned Fixes
          </span>
          {data?.diff && (
            <button
              onClick={copyDiff}
              className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-200 text-xs font-medium border border-slate-700 transition"
            >
              {copied ? <Check className="w-3.5 h-3.5 text-emerald-400" /> : <Copy className="w-3.5 h-3.5" />}
              {copied ? 'Patch Copied' : 'Copy Unified Diff'}
            </button>
          )}
        </div>
      </div>

      <div className="p-4 space-y-4">
        {fixes.length === 0 ? (
          <div className="text-center py-6 text-slate-500 text-xs">
            No active compliance violations requiring automated remediation.
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
            {fixes.map((f, i) => (
              <div key={i} className="p-3 rounded-lg bg-slate-900 border border-slate-800 text-xs space-y-1">
                <div className="flex items-center justify-between">
                  <span className="font-semibold text-white flex items-center gap-1.5">
                    <Sparkles className="w-3.5 h-3.5 text-blue-400" />
                    {f.type}
                  </span>
                  <span className="font-mono text-slate-400 text-[11px]">
                    {f.file_path}:{f.line}
                  </span>
                </div>
                <p className="text-slate-400 text-[11px]">{f.description}</p>
              </div>
            ))}
          </div>
        )}

        {data?.diff && (
          <div className="space-y-1.5">
            <h4 className="text-xs font-bold text-slate-400 uppercase tracking-wider">
              Unified Git Patch Preview
            </h4>
            <pre className="p-4 rounded-lg bg-slate-950 border border-slate-800 text-[11px] font-mono overflow-x-auto leading-relaxed">
              {data.diff.split('\n').map((line, idx) => {
                let colorClass = 'text-slate-400';
                if (line.startsWith('+') && !line.startsWith('+++')) {
                  colorClass = 'text-emerald-400 bg-emerald-950/20';
                } else if (line.startsWith('-') && !line.startsWith('---')) {
                  colorClass = 'text-rose-400 bg-rose-950/20';
                } else if (line.startsWith('@@') || line.startsWith('---') || line.startsWith('+++')) {
                  colorClass = 'text-sky-400 font-bold';
                }
                return (
                  <div key={idx} className={colorClass}>
                    {line}
                  </div>
                );
              })}
            </pre>
          </div>
        )}
      </div>
    </div>
  );
};
