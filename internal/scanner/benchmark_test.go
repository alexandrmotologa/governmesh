package scanner

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/alexandrmotologa/governmesh/internal/config"
)

func BenchmarkScannerMemory(b *testing.B) {
	sampleDir, err := filepath.Abs("../../examples/enterprise_ai_service")
	if err != nil {
		b.Fatal(err)
	}

	cfg := config.DefaultConfig()
	scn := NewScanner(cfg, nil)

	var memStatsBefore, memStatsAfter runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&memStatsBefore)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := scn.Run(sampleDir)
		if err != nil {
			b.Fatal(err)
		}
	}
	b.StopTimer()

	runtime.ReadMemStats(&memStatsAfter)
	allocatedMB := float64(memStatsAfter.Alloc) / (1024 * 1024)

	b.Logf("Total Heap Alloc: %.2f MB (Target < 20MB)", allocatedMB)
	if allocatedMB > 20.0 {
		b.Errorf("Memory exceeded 20MB threshold: %.2f MB", allocatedMB)
	}
}
