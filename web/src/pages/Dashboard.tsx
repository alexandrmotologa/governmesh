import React, { useState } from 'react';
import { ScanResult } from '../types';
import { RiskMatrixChart } from '../components/RiskMatrixChart';
import { ComplianceScore } from '../components/ComplianceScore';
import { AibomExplorer } from '../components/AibomExplorer';
import { PiiEgressTable } from '../components/PiiEgressTable';
import { AttestationCard } from '../components/AttestationCard';
import { LayoutDashboard, Database, ShieldAlert, CheckSquare } from 'lucide-react';

interface DashboardProps {
  data: ScanResult;
}

export const Dashboard: React.FC<DashboardProps> = ({ data }) => {
  const [activeTab, setActiveTab] = useState<'overview' | 'aibom' | 'pii' | 'standards'>('overview');

  const tabs = [
    { id: 'overview', label: 'Executive Overview', icon: LayoutDashboard },
    { id: 'aibom', label: `AIBOM Inventory (${data.assets.length})`, icon: Database },
    { id: 'pii', label: `PII & Data Egress (${data.pii_exposures.length})`, icon: ShieldAlert },
    { id: 'standards', label: 'ISO 42001 & NIST RMF', icon: CheckSquare },
  ];

  return (
    <div className="max-w-7xl mx-auto px-6 py-6 space-y-6">
      {/* Navigation tabs */}
      <div className="flex border-b border-slate-800 gap-1 overflow-x-auto">
        {tabs.map((tab) => {
          const Icon = tab.icon;
          const isActive = activeTab === tab.id;
          return (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id as any)}
              className={`flex items-center gap-2 px-4 py-2.5 text-xs font-semibold rounded-t-lg transition border-b-2 whitespace-nowrap ${
                isActive
                  ? 'border-blue-500 text-white bg-slate-900/50'
                  : 'border-transparent text-slate-400 hover:text-slate-200 hover:bg-slate-900/20'
              }`}
            >
              <Icon className="w-4 h-4" />
              {tab.label}
            </button>
          );
        })}
      </div>

      {/* Tab Content */}
      {activeTab === 'overview' && (
        <div className="space-y-6">
          <RiskMatrixChart data={data} />
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <ComplianceScore data={data} />
            <AttestationCard data={data} />
          </div>
          <AibomExplorer assets={data.assets} />
        </div>
      )}

      {activeTab === 'aibom' && (
        <div className="space-y-6">
          <AibomExplorer assets={data.assets} />
        </div>
      )}

      {activeTab === 'pii' && (
        <div className="space-y-6">
          <PiiEgressTable exposures={data.pii_exposures} />
        </div>
      )}

      {activeTab === 'standards' && (
        <div className="bg-[#0f172a] rounded-xl border border-slate-800 p-6 space-y-6">
          <div>
            <h3 className="text-base font-bold text-white">ISO/IEC 42001 & NIST AI RMF Control Mapping</h3>
            <p className="text-xs text-slate-400 mt-1 leading-relaxed">
              Standardized compliance checklist mapping static AST evidence and technical safeguards to mandatory controls.
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            {/* ISO 42001 */}
            <div className="space-y-3">
              <h4 className="text-sm font-semibold text-blue-400 border-b border-slate-800 pb-2">
                ISO/IEC 42001:2023 Controls
              </h4>
              <div className="space-y-2 text-xs">
                <div className="p-3 rounded bg-slate-900/80 border border-slate-800/80">
                  <div className="font-semibold text-white">A.5.1 AI Policy & Governance</div>
                  <p className="text-slate-400 mt-0.5">Automated quality gate thresholds active for risk tiering.</p>
                </div>
                <div className="p-3 rounded bg-slate-900/80 border border-slate-800/80">
                  <div className="font-semibold text-white">A.6.1 AI Risk Assessment</div>
                  <p className="text-slate-400 mt-0.5">Continuous evaluation against EU AI Act Annex III high-risk domains.</p>
                </div>
                <div className="p-3 rounded bg-slate-900/80 border border-slate-800/80">
                  <div className="font-semibold text-white">A.7.2 Data for AI Systems & Privacy</div>
                  <p className="text-slate-400 mt-0.5">Static prompt inspection for unmasked PII/PHI variable egress.</p>
                </div>
                <div className="p-3 rounded bg-slate-900/80 border border-slate-800/80">
                  <div className="font-semibold text-white">A.8.2 AI System Inventory & Traceability</div>
                  <p className="text-slate-400 mt-0.5">CycloneDX 1.6 AIBOM with SHA-256 Merkle chain attestation.</p>
                </div>
              </div>
            </div>

            {/* NIST AI RMF */}
            <div className="space-y-3">
              <h4 className="text-sm font-semibold text-purple-400 border-b border-slate-800 pb-2">
                NIST AI Risk Management Framework (AI RMF 1.0)
              </h4>
              <div className="space-y-2 text-xs">
                <div className="p-3 rounded bg-slate-900/80 border border-slate-800/80">
                  <div className="font-semibold text-white">GOVERN 1.1 / 2.1</div>
                  <p className="text-slate-400 mt-0.5">Legal & regulatory requirements mapped; authoritative asset inventory maintained.</p>
                </div>
                <div className="p-3 rounded bg-slate-900/80 border border-slate-800/80">
                  <div className="font-semibold text-white">MAP 1.1 / 1.5</div>
                  <p className="text-slate-400 mt-0.5">Model invocation scope and high-risk domain identifiers categorized.</p>
                </div>
                <div className="p-3 rounded bg-slate-900/80 border border-slate-800/80">
                  <div className="font-semibold text-white">MEASURE 1.1</div>
                  <p className="text-slate-400 mt-0.5">Perimeter audits for sensitive demographic and personal data egress.</p>
                </div>
                <div className="p-3 rounded bg-slate-900/80 border border-slate-800/80">
                  <div className="font-semibold text-white">MANAGE 2.4</div>
                  <p className="text-slate-400 mt-0.5">Automated quality gate blocking non-compliant builds in CI/CD.</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};
