package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/alexandrmotologa/governmesh/internal/aibom"
	"github.com/alexandrmotologa/governmesh/internal/attestation"
	"github.com/alexandrmotologa/governmesh/internal/classifier"
	"github.com/alexandrmotologa/governmesh/internal/config"
	"github.com/alexandrmotologa/governmesh/internal/report"
	"github.com/alexandrmotologa/governmesh/internal/scanner"
)

// Server coordinates the local Governance Studio API and embedded frontend.
type Server struct {
	cfg         *config.Config
	targetDir   string
	mu          sync.RWMutex
	lastResult  *scanner.ScanResult
	lastAIBOM   *aibom.CycloneDXBOM
	lastSPDX    *aibom.SPDXDocument
	lastAttest  *attestation.Manifest
	sseClients  map[chan string]bool
	sseMu       sync.Mutex
}

// NewServer initializes a new Server.
func NewServer(cfg *config.Config, targetDir string) *Server {
	if cfg == nil {
		cfg = config.DefaultConfig()
	}
	return &Server{
		cfg:        cfg,
		targetDir:  targetDir,
		sseClients: make(map[chan string]bool),
	}
}

// RegisterRoutes sets up all REST and SSE endpoints.
func (s *Server) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/health", s.handleHealth)
	mux.HandleFunc("/api/scan", s.handleScan)
	mux.HandleFunc("/api/scan/stream", s.handleScanStream)
	mux.HandleFunc("/api/aibom", s.handleAIBOM)
	mux.HandleFunc("/api/spdx", s.handleSPDX)
	mux.HandleFunc("/api/attestation", s.handleAttestation)
	mux.HandleFunc("/api/report/html", s.handleReportHTML)
	mux.HandleFunc("/api/report/md", s.handleReportMD)
	mux.HandleFunc("/api/verify", s.handleVerify)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "healthy",
		"service": "governmesh-studio",
		"version": "1.0.0",
	})
}

func (s *Server) handleScanStream(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	messageChan := make(chan string, 10)
	s.sseMu.Lock()
	s.sseClients[messageChan] = true
	s.sseMu.Unlock()

	defer func() {
		s.sseMu.Lock()
		delete(s.sseClients, messageChan)
		close(messageChan)
		s.sseMu.Unlock()
	}()

	notify := r.Context().Done()
	for {
		select {
		case <-notify:
			return
		case msg := <-messageChan:
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()
		}
	}
}

func (s *Server) broadcastSSE(msg string) {
	s.sseMu.Lock()
	defer s.sseMu.Unlock()
	for ch := range s.sseClients {
		select {
		case ch <- msg:
		default:
		}
	}
}

func (s *Server) handleScan(w http.ResponseWriter, r *http.Request) {
	target := s.targetDir
	if r.Method == http.MethodPost {
		var req struct {
			Path string `json:"path"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err == nil && req.Path != "" {
			target = req.Path
		}
	}

	s.broadcastSSE(fmt.Sprintf(`{"type":"started","target":"%s"}`, target))

	scn := scanner.NewScanner(s.cfg, func(curFile string, files int, assets int) {
		s.broadcastSSE(fmt.Sprintf(`{"type":"progress","file":"%s","scanned":%d,"assets":%d}`, curFile, files, assets))
	})

	res, err := scn.Run(target)
	if err != nil {
		s.broadcastSSE(fmt.Sprintf(`{"type":"error","message":"%s"}`, err.Error()))
		http.Error(w, fmt.Sprintf("scan failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Classify
	classifier.ProcessScanResult(res, s.cfg)

	// AIBOM & Attestation
	cdx, _ := aibom.GenerateCycloneDX(res)
	spdxDoc, _ := aibom.GenerateSPDX(res)
	cdxBytes, _ := cdx.ToJSON()
	att, _ := attestation.GenerateManifest(res, cdxBytes, s.cfg)

	s.mu.Lock()
	s.lastResult = res
	s.lastAIBOM = cdx
	s.lastSPDX = spdxDoc
	s.lastAttest = att
	s.mu.Unlock()

	s.broadcastSSE(`{"type":"completed"}`)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (s *Server) handleAIBOM(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.lastAIBOM == nil {
		http.Error(w, "No scan result available. Run a scan first.", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	data, _ := s.lastAIBOM.ToJSON()
	w.Write(data)
}

func (s *Server) handleSPDX(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.lastSPDX == nil {
		http.Error(w, "No scan result available. Run a scan first.", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	data, _ := s.lastSPDX.ToJSON()
	w.Write(data)
}

func (s *Server) handleAttestation(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.lastAttest == nil {
		http.Error(w, "No attestation available.", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.lastAttest)
}

func (s *Server) handleReportHTML(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	res := s.lastResult
	att := s.lastAttest
	s.mu.RUnlock()

	if res == nil {
		http.Error(w, "No scan result available.", http.StatusNotFound)
		return
	}

	htmlContent := report.GenerateHTMLReport(res, att)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(htmlContent))
}

func (s *Server) handleReportMD(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	res := s.lastResult
	att := s.lastAttest
	s.mu.RUnlock()

	if res == nil {
		http.Error(w, "No scan result available.", http.StatusNotFound)
		return
	}

	mdContent := report.GenerateMarkdownReport(res, att)
	w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
	w.Write([]byte(mdContent))
}

func (s *Server) handleVerify(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST required", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "reading body: "+err.Error(), http.StatusBadRequest)
		return
	}

	var req struct {
		ManifestPath string `json:"manifest_path"`
		AIBOMPath    string `json:"aibom_path"`
		TargetDir    string `json:"target_dir"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "invalid json payload", http.StatusBadRequest)
		return
	}

	result, err := attestation.VerifyBundle(req.ManifestPath, req.AIBOMPath, req.TargetDir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
