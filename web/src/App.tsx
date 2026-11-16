import React, { useState, useEffect } from 'react';
import { ScanResult } from './types';
import { Header } from './components/Header';
import { Dashboard } from './pages/Dashboard';
import { ShieldCheck, AlertCircle, RefreshCw } from 'lucide-react';

export const App: React.FC = () => {
  const [data, setData] = useState<ScanResult | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [liveLog, setLiveLog] = useState<string | null>(null);

  const triggerScan = async () => {
    setLoading(true);
    setError(null);
    setLiveLog('Initializing repository scan...');
    try {
      const res = await fetch('/api/scan', { method: 'POST' });
      if (!res.ok) {
        throw new Error(`Scan request failed with HTTP ${res.status}`);
      }
      const json: ScanResult = await res.json();
      setData(json);
      setLiveLog(null);
    } catch (err: any) {
      setError(err.message || 'Error occurred during scan execution');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    // Initial scan on load
    triggerScan();

    // Setup Server-Sent Events (SSE) for live streaming progress
    const sse = new EventSource('/api/scan/stream');
    sse.onmessage = (event) => {
      try {
        const msg = JSON.parse(event.data);
        if (msg.type === 'progress') {
          setLiveLog(`Scanning: ${msg.file} (${msg.scanned} files inspected, ${msg.assets} AI assets found)`);
        } else if (msg.type === 'completed') {
          setLiveLog(null);
        }
      } catch (e) {
        // Fallback for non-json SSE
      }
    };

    return () => {
      sse.close();
    };
  }, []);

  return (
    <div className="min-h-screen bg-[#090d16] text-slate-100 flex flex-col font-sans">
      <Header data={data} loading={loading} onTriggerScan={triggerScan} />

      {liveLog && (
        <div className="bg-blue-950/60 border-b border-blue-900/50 px-6 py-2 text-xs text-blue-300 flex items-center justify-between">
          <div className="flex items-center gap-2 truncate">
            <RefreshCw className="w-3.5 h-3.5 animate-spin flex-shrink-0" />
            <span className="truncate">{liveLog}</span>
          </div>
        </div>
      )}

      {error && (
        <div className="max-w-7xl mx-auto px-6 pt-6 w-full">
          <div className="p-4 rounded-xl bg-rose-950/50 border border-rose-900 text-rose-300 text-xs flex items-center gap-2">
            <AlertCircle className="w-4 h-4 flex-shrink-0 text-rose-400" />
            <span>{error}</span>
          </div>
        </div>
      )}

      <main className="flex-1">
        {loading && !data ? (
          <div className="flex flex-col items-center justify-center min-h-[60vh] gap-3">
            <RefreshCw className="w-8 h-8 text-blue-500 animate-spin" />
            <h2 className="text-sm font-semibold text-slate-200">Executing Static AST AI Codebase Audit...</h2>
            <p className="text-xs text-slate-400">Inspecting Python, TypeScript, Go, manifests and prompt templates.</p>
          </div>
        ) : data ? (
          <Dashboard data={data} />
        ) : null}
      </main>

      <footer className="border-t border-slate-900 bg-[#070b13] px-6 py-4 text-center text-xs text-slate-500">
        governmesh v1.0.0 &bull; Automated EU AI Act, NIST AI RMF & ISO 42001 Compliance Architecture
      </footer>
    </div>
  );
};

export default App;
