import React, { useState, useEffect } from 'react';
import { BookOpen, Download, Copy, Check } from 'lucide-react';

export const AnnexIVViewer: React.FC = () => {
  const [content, setContent] = useState<string>('');
  const [loading, setLoading] = useState<boolean>(true);
  const [copied, setCopied] = useState<boolean>(false);

  useEffect(() => {
    fetch('/api/annex-iv')
      .then((res) => res.text())
      .then((text) => {
        setContent(text);
        setLoading(false);
      })
      .catch(() => setLoading(false));
  }, []);

  const copyDossier = () => {
    navigator.clipboard.writeText(content);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  const downloadDossier = () => {
    const blob = new Blob([content], { type: 'text/markdown' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'eu_ai_act_annex_iv_technical_dossier.md';
    a.click();
    URL.revokeObjectURL(url);
  };

  return (
    <div className="bg-[#0f172a] rounded-xl border border-slate-800 overflow-hidden space-y-4">
      <div className="p-4 border-b border-slate-800 flex justify-between items-center">
        <div>
          <h3 className="text-sm font-semibold text-white flex items-center gap-2">
            <BookOpen className="w-4 h-4 text-blue-400" />
            EU AI Act Annex IV Statutory Technical Documentation
          </h3>
          <p className="text-xs text-slate-400 mt-0.5">
            Mandatory conformity dossier required for High-Risk AI systems under Article 11 and Annex IV of Regulation (EU) 2024/1689.
          </p>
        </div>
        <div className="flex items-center gap-2">
          <button
            onClick={copyDossier}
            className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-200 text-xs font-medium border border-slate-700 transition"
          >
            {copied ? <Check className="w-3.5 h-3.5 text-emerald-400" /> : <Copy className="w-3.5 h-3.5" />}
            {copied ? 'Copied' : 'Copy Markdown'}
          </button>
          <button
            onClick={downloadDossier}
            className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-blue-600 hover:bg-blue-500 text-white text-xs font-medium transition"
          >
            <Download className="w-3.5 h-3.5" />
            Download Dossier (.md)
          </button>
        </div>
      </div>

      <div className="p-4">
        {loading ? (
          <div className="text-center py-8 text-slate-400 text-xs">
            Compiling technical documentation from AST evidence...
          </div>
        ) : (
          <pre className="p-4 rounded-lg bg-slate-950 border border-slate-800 text-[11px] font-mono text-slate-300 overflow-x-auto whitespace-pre-wrap leading-relaxed">
            {content}
          </pre>
        )}
      </div>
    </div>
  );
};
