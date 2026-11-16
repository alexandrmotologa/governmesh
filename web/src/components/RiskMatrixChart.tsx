import React from 'react';
import { AlertOctagon, AlertTriangle, Eye, CheckCircle2 } from 'lucide-react';
import { ScanResult } from '../types';

interface RiskMatrixChartProps {
  data: ScanResult;
}

export const RiskMatrixChart: React.FC<RiskMatrixChartProps> = ({ data }) => {
  const counts = data.risk_summary;
  const total = data.assets.length || 1;

  const tiers = [
    {
      id: 'unacceptable',
      label: 'Unacceptable Risk',
      article: 'EU AI Act Article 5',
      count: counts.unacceptable || 0,
      icon: AlertOctagon,
      textColor: 'text-rose-400',
      bgColor: 'bg-rose-500/10',
      borderColor: 'border-rose-500/30',
      desc: 'Prohibited practices (manipulation, social scoring, unauthorized biometrics).',
    },
    {
      id: 'high',
      label: 'High-Risk',
      article: 'Annex III / Article 6',
      count: counts.high || 0,
      icon: AlertTriangle,
      textColor: 'text-amber-400',
      bgColor: 'bg-amber-500/10',
      borderColor: 'border-amber-500/30',
      desc: 'Critical infrastructure, HR recruitment, credit scoring, medical decisions.',
    },
    {
      id: 'transparency',
      label: 'Specific Transparency',
      article: 'EU AI Act Article 50',
      count: counts.transparency || 0,
      icon: Eye,
      textColor: 'text-sky-400',
      bgColor: 'bg-sky-500/10',
      borderColor: 'border-sky-500/30',
      desc: 'Chatbots, generative models, and synthetic content requiring disclosure.',
    },
    {
      id: 'minimal',
      label: 'Minimal / Low Risk',
      article: 'Article 69 / Recital 68',
      count: counts.minimal || 0,
      icon: CheckCircle2,
      textColor: 'text-emerald-400',
      bgColor: 'bg-emerald-500/10',
      borderColor: 'border-emerald-500/30',
      desc: 'Internal developer tooling, search indexing, translation, spam filters.',
    },
  ];

  return (
    <div className="space-y-4">
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {tiers.map((tier) => {
          const Icon = tier.icon;
          return (
            <div
              key={tier.id}
              className={`p-4 rounded-xl border ${tier.bgColor} ${tier.borderColor} flex flex-col justify-between`}
            >
              <div>
                <div className="flex items-center justify-between mb-2">
                  <span className="text-xs font-semibold text-slate-400 tracking-wider uppercase">
                    {tier.article}
                  </span>
                  <Icon className={`w-5 h-5 ${tier.textColor}`} />
                </div>
                <div className="text-2xl font-black text-white">{tier.count}</div>
                <div className={`text-sm font-semibold mt-0.5 ${tier.textColor}`}>{tier.label}</div>
              </div>
              <p className="text-xs text-slate-400 mt-3 border-t border-slate-800/80 pt-2 leading-relaxed">
                {tier.desc}
              </p>
            </div>
          );
        })}
      </div>

      {/* Distribution progress bar */}
      <div className="bg-[#0f172a] p-4 rounded-xl border border-slate-800">
        <div className="flex justify-between items-center mb-2">
          <span className="text-xs font-semibold text-slate-400 uppercase tracking-wider">
            Risk Tier Distribution
          </span>
          <span className="text-xs text-slate-500">{data.assets.length} AI systems analyzed</span>
        </div>
        <div className="w-full h-3 rounded-full bg-slate-800 overflow-hidden flex">
          {counts.unacceptable > 0 && (
            <div
              style={{ width: `${(counts.unacceptable / total) * 100}%` }}
              className="bg-rose-500 h-full"
              title={`Unacceptable: ${counts.unacceptable}`}
            />
          )}
          {counts.high > 0 && (
            <div
              style={{ width: `${(counts.high / total) * 100}%` }}
              className="bg-amber-500 h-full"
              title={`High-Risk: ${counts.high}`}
            />
          )}
          {counts.transparency > 0 && (
            <div
              style={{ width: `${(counts.transparency / total) * 100}%` }}
              className="bg-sky-500 h-full"
              title={`Transparency: ${counts.transparency}`}
            />
          )}
          {counts.minimal > 0 && (
            <div
              style={{ width: `${(counts.minimal / total) * 100}%` }}
              className="bg-emerald-500 h-full"
              title={`Minimal: ${counts.minimal}`}
            />
          )}
        </div>
      </div>
    </div>
  );
};
