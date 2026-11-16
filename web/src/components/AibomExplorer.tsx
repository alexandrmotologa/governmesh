import React, { useState } from 'react';
import { Search, Filter, ChevronDown, ChevronRight, Cpu, Cloud, Database, ShieldAlert } from 'lucide-react';
import { DiscoveredAsset, RiskLevel } from '../types';

interface AibomExplorerProps {
  assets: DiscoveredAsset[];
}

export const AibomExplorer: React.FC<AibomExplorerProps> = ({ assets }) => {
  const [search, setSearch] = useState('');
  const [filterTier, setFilterTier] = useState<string>('all');
  const [expandedId, setExpandedId] = useState<string | null>(null);

  const filtered = assets.filter((a) => {
    const matchesSearch =
      a.model_family.toLowerCase().includes(search.toLowerCase()) ||
      a.source_file_path.toLowerCase().includes(search.toLowerCase()) ||
      a.provider.toLowerCase().includes(search.toLowerCase());
    const matchesTier = filterTier === 'all' || a.risk.tier === filterTier;
    return matchesSearch && matchesTier;
  });

  const getTierBadge = (tier: RiskLevel) => {
    switch (tier) {
      case 'unacceptable':
        return <span className="px-2 py-0.5 rounded text-xs font-bold bg-rose-500/20 text-rose-400 border border-rose-500/30">Unacceptable</span>;
      case 'high':
        return <span className="px-2 py-0.5 rounded text-xs font-bold bg-amber-500/20 text-amber-400 border border-amber-500/30">High-Risk</span>;
      case 'transparency':
        return <span className="px-2 py-0.5 rounded text-xs font-bold bg-sky-500/20 text-sky-400 border border-sky-500/30">Transparency</span>;
      default:
        return <span className="px-2 py-0.5 rounded text-xs font-bold bg-emerald-500/20 text-emerald-400 border border-emerald-500/30">Minimal</span>;
    }
  };

  return (
    <div className="bg-[#0f172a] rounded-xl border border-slate-800 overflow-hidden">
      <div className="p-4 border-b border-slate-800 flex flex-col sm:flex-row gap-3 justify-between items-center">
        <div className="relative w-full sm:w-72">
          <Search className="w-4 h-4 absolute left-3 top-3 text-slate-500" />
          <input
            type="text"
            placeholder="Search models, files, providers..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="w-full pl-9 pr-3 py-2 bg-slate-900 border border-slate-700/80 rounded-lg text-xs text-white placeholder-slate-500 focus:outline-none focus:border-blue-500"
          />
        </div>

        <div className="flex items-center gap-2 w-full sm:w-auto">
          <Filter className="w-4 h-4 text-slate-500" />
          <select
            value={filterTier}
            onChange={(e) => setFilterTier(e.target.value)}
            className="bg-slate-900 border border-slate-700/80 rounded-lg px-3 py-2 text-xs text-slate-300 focus:outline-none focus:border-blue-500"
          >
            <option value="all">All Risk Tiers ({assets.length})</option>
            <option value="unacceptable">Unacceptable Risk</option>
            <option value="high">High-Risk (Annex III)</option>
            <option value="transparency">Specific Transparency</option>
            <option value="minimal">Minimal Risk</option>
          </select>
        </div>
      </div>

      <div className="overflow-x-auto">
        <table className="w-full text-left border-collapse">
          <thead>
            <tr className="bg-slate-900/60 text-slate-400 text-xs font-semibold uppercase tracking-wider border-b border-slate-800">
              <th className="py-3 px-4 w-10"></th>
              <th className="py-3 px-4">AI Model / Component</th>
              <th className="py-3 px-4">Provider</th>
              <th className="py-3 px-4">Risk Tier</th>
              <th className="py-3 px-4">Framework</th>
              <th className="py-3 px-4">Egress Target</th>
              <th className="py-3 px-4">Source Location</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-800 text-xs">
            {filtered.length === 0 ? (
              <tr>
                <td colSpan={7} className="py-8 text-center text-slate-500">
                  No AI components match the selected query.
                </td>
              </tr>
            ) : (
              filtered.map((asset) => {
                const isExpanded = expandedId === asset.id;
                return (
                  <React.Fragment key={asset.id}>
                    <tr
                      onClick={() => setExpandedId(isExpanded ? null : asset.id)}
                      className="hover:bg-slate-800/40 cursor-pointer transition"
                    >
                      <td className="py-3 px-4 text-slate-500">
                        {isExpanded ? <ChevronDown className="w-4 h-4" /> : <ChevronRight className="w-4 h-4" />}
                      </td>
                      <td className="py-3 px-4 font-semibold text-white flex items-center gap-2">
                        <Cpu className="w-3.5 h-3.5 text-blue-400 flex-shrink-0" />
                        {asset.model_family}
                      </td>
                      <td className="py-3 px-4 text-slate-300">{asset.provider}</td>
                      <td className="py-3 px-4">{getTierBadge(asset.risk.tier)}</td>
                      <td className="py-3 px-4">
                        <span className="px-2 py-0.5 rounded bg-slate-800 text-slate-300 font-mono text-[11px]">
                          {asset.framework}
                        </span>
                      </td>
                      <td className="py-3 px-4">
                        <span
                          className={`font-mono text-[11px] ${
                            asset.egress_target === 'external_cloud'
                              ? 'text-sky-400'
                              : 'text-emerald-400'
                          }`}
                        >
                          {asset.egress_target}
                        </span>
                      </td>
                      <td className="py-3 px-4 text-slate-400 font-mono text-[11px]">
                        {asset.source_file_path}:{asset.line_number}
                      </td>
                    </tr>

                    {isExpanded && (
                      <tr className="bg-slate-900/40">
                        <td colSpan={7} className="p-4 border-t border-slate-800/80">
                          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                            <div className="space-y-3">
                              <div>
                                <h4 className="text-[11px] font-bold text-slate-400 uppercase tracking-wider">
                                  Regulatory Classification Details
                                </h4>
                                <div className="mt-1 text-xs text-slate-200">
                                  <strong>Article:</strong> {asset.risk.applicable_article}
                                </div>
                                <p className="mt-1 text-xs text-slate-400 leading-relaxed">
                                  {asset.risk.justification}
                                </p>
                              </div>

                              {asset.risk.prohibited_reason && (
                                <div className="p-2.5 rounded bg-rose-950/40 border border-rose-900 text-rose-300 text-xs">
                                  <strong>Prohibition Note:</strong> {asset.risk.prohibited_reason}
                                </div>
                              )}

                              <div>
                                <h4 className="text-[11px] font-bold text-slate-400 uppercase tracking-wider">
                                  Standards Controls Mapped
                                </h4>
                                <div className="flex flex-wrap gap-1.5 mt-1.5">
                                  {asset.risk.nist_categories?.map((n) => (
                                    <span key={n} className="px-2 py-0.5 bg-blue-950/50 text-blue-300 border border-blue-800/50 rounded text-[11px]">
                                      {n}
                                    </span>
                                  ))}
                                  {asset.risk.iso_42001_controls?.map((iso) => (
                                    <span key={iso} className="px-2 py-0.5 bg-purple-950/50 text-purple-300 border border-purple-800/50 rounded text-[11px]">
                                      {iso}
                                    </span>
                                  ))}
                                </div>
                              </div>
                            </div>

                            <div className="space-y-3">
                              <div>
                                <h4 className="text-[11px] font-bold text-slate-400 uppercase tracking-wider">
                                  Code Context & Invocations
                                </h4>
                                <pre className="mt-1 p-3 rounded bg-slate-950 border border-slate-800 text-slate-300 font-mono text-[11px] overflow-x-auto">
                                  {asset.context_snippet || '// AST context snippet unavailable'}
                                </pre>
                              </div>

                              {asset.prompt_vars && asset.prompt_vars.length > 0 && (
                                <div>
                                  <h4 className="text-[11px] font-bold text-slate-400 uppercase tracking-wider">
                                    Interpolated Prompt Variables
                                  </h4>
                                  <div className="flex flex-wrap gap-1 mt-1 font-mono text-[11px] text-amber-300">
                                    {asset.prompt_vars.map((v) => (
                                      <span key={v} className="px-1.5 py-0.5 bg-amber-950/30 rounded border border-amber-900/40">
                                        {`{${v}}`}
                                      </span>
                                    ))}
                                  </div>
                                </div>
                              )}
                            </div>
                          </div>
                        </td>
                      </tr>
                    )}
                  </React.Fragment>
                );
              })
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
};
