package scanner

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/alexandrmotologa/governmesh/internal/config"
)

// ProgressCallback provides real-time notifications during scanning.
type ProgressCallback func(currentFile string, filesScanned int, assetsFound int)

// Scanner coordinates repository traversal and multi-language AST inspection.
type Scanner struct {
	cfg      *config.Config
	callback ProgressCallback
}

// NewScanner creates a new Scanner instance.
func NewScanner(cfg *config.Config, cb ProgressCallback) *Scanner {
	if cfg == nil {
		cfg = config.DefaultConfig()
	}
	return &Scanner{
		cfg:      cfg,
		callback: cb,
	}
}

// Run executes the full scan against the target root directory.
func (s *Scanner) Run(rootDir string) (*ScanResult, error) {
	startTime := time.Now()
	absRoot, err := filepath.Abs(rootDir)
	if err != nil {
		return nil, fmt.Errorf("resolving absolute path: %w", err)
	}

	info, err := os.Stat(absRoot)
	if err != nil {
		return nil, fmt.Errorf("stat target directory: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("target path is not a directory: %s", absRoot)
	}

	totalFiles := 0
	scannedAIFiles := 0
	var allAssets []DiscoveredAsset

	excludeMap := make(map[string]bool)
	for _, excl := range s.cfg.Scan.ExcludePaths {
		excludeMap[excl] = true
	}

	err = filepath.WalkDir(absRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // Skip unreadable paths
		}

		// Check directory exclusions
		if d.IsDir() {
			base := d.Name()
			if excludeMap[base] || strings.HasPrefix(base, ".") && base != "." {
				return filepath.SkipDir
			}
			return nil
		}

		totalFiles++
		base := d.Name()
		ext := strings.ToLower(filepath.Ext(path))

		// Check file size limit
		if s.cfg.Scan.MaxFileSizeMB > 0 {
			finfo, err := d.Info()
			if err == nil && finfo.Size() > (s.cfg.Scan.MaxFileSizeMB*1024*1024) {
				return nil
			}
		}

		var assets []DiscoveredAsset
		var scanErr error

		switch {
		case ext == ".py":
			assets, scanErr = ScanPythonFile(path, absRoot)
		case ext == ".ts" || ext == ".tsx" || ext == ".js" || ext == ".jsx":
			assets, scanErr = ScanTypeScriptFile(path, absRoot)
		case ext == ".go":
			assets, scanErr = ScanGoFile(path, absRoot)
		case ext == ".prompt" || (strings.Contains(strings.ToLower(base), "prompt") && (ext == ".yaml" || ext == ".yml" || ext == ".json" || ext == ".txt")):
			assets, scanErr = ScanPromptFile(path, absRoot)
		case strings.Contains(strings.ToLower(base), "dockerfile") || ext == ".env" || ext == ".yaml" || ext == ".yml":
			if s.cfg.Scan.ScanManifests {
				assets, scanErr = ScanManifestFile(path, absRoot)
			}
		}

		if scanErr == nil && len(assets) > 0 {
			scannedAIFiles++
			allAssets = append(allAssets, assets...)
		}

		if s.callback != nil && totalFiles%20 == 0 {
			s.callback(path, totalFiles, len(allAssets))
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("walking directory tree: %w", err)
	}

	// Calculate Merkle hash chain
	sort.Slice(allAssets, func(i, j int) bool {
		return allAssets[i].ID < allAssets[j].ID
	})

	hasher := sha256.New()
	for _, a := range allAssets {
		hasher.Write([]byte(a.ID + a.FileHash + a.ModelFamily))
	}
	hashChain := hex.EncodeToString(hasher.Sum(nil))

	duration := time.Since(startTime).Milliseconds()

	result := &ScanResult{
		ProjectName:    s.cfg.ProjectName,
		ScannedPath:    absRoot,
		ScanTimestamp:  startTime,
		DurationMs:     duration,
		TotalFiles:     totalFiles,
		ScannedAIFiles: scannedAIFiles,
		Assets:         allAssets,
		RiskSummary:    make(map[config.RiskLevel]int),
		HashChain:      hashChain,
	}

	return result, nil
}
