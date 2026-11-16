package scanner

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Microservice holds information about a recognized sub-project in a monorepo.
type Microservice struct {
	Name       string   `json:"name"`
	Path       string   `json:"path"`
	TechStack  string   `json:"tech_stack"` // "python", "typescript", "go", "container"
	AssetIDs   []string `json:"asset_ids"`
}

// MonorepoInspection represents cross-service dependency discoveries.
type MonorepoInspection struct {
	IsMonorepo    bool           `json:"is_monorepo"`
	Services      []Microservice `json:"services"`
	CrossEgresses []string       `json:"cross_egresses"`
}

var serviceMarkerFiles = map[string]string{
	"package.json":     "typescript/node",
	"pyproject.toml":   "python",
	"requirements.txt": "python",
	"go.mod":           "go",
	"Dockerfile":       "container",
}

// DetectMonorepo scans the directory structure for microservice boundaries.
func DetectMonorepo(rootDir string, assets []DiscoveredAsset) *MonorepoInspection {
	servicesMap := make(map[string]*Microservice)

	filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || !d.IsDir() {
			return nil
		}

		base := d.Name()
		if base == ".git" || base == "node_modules" || base == "vendor" || base == ".venv" {
			return filepath.SkipDir
		}

		// Check if this directory contains a package marker
		entries, err := os.ReadDir(path)
		if err != nil {
			return nil
		}

		for _, e := range entries {
			if tech, exists := serviceMarkerFiles[e.Name()]; exists && path != rootDir {
				relPath, _ := filepath.Rel(rootDir, path)
				relPath = filepath.ToSlash(relPath)
				if _, already := servicesMap[relPath]; !already {
					serviceName := filepath.Base(path)
					servicesMap[relPath] = &Microservice{
						Name:      serviceName,
						Path:      relPath,
						TechStack: tech,
					}
				}
			}
		}

		return nil
	})

	// Associate discovered assets with microservices
	for _, a := range assets {
		for svcPath, svc := range servicesMap {
			if strings.HasPrefix(a.SourceFilePath, svcPath+"/") || a.SourceFilePath == svcPath {
				svc.AssetIDs = append(svc.AssetIDs, a.ID)
			}
		}
	}

	var servicesList []Microservice
	for _, svc := range servicesMap {
		servicesList = append(servicesList, *svc)
	}

	isMonorepo := len(servicesList) > 1

	return &MonorepoInspection{
		IsMonorepo: isMonorepo,
		Services:   servicesList,
	}
}
