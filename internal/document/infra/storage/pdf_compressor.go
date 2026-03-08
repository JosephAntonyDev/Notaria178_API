package storage

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

func CompressPDF(filePath string) {
	if !strings.EqualFold(filepath.Ext(filePath), ".pdf") {
		return
	}

	originalInfo, err := os.Stat(filePath)
	if err != nil {
		return
	}

	tmpPath := filePath + ".optimized"

	if err := api.OptimizeFile(filePath, tmpPath, nil); err != nil {
		os.Remove(tmpPath)
		log.Printf("[PDF] Optimización omitida para %s: %v", filepath.Base(filePath), err)
		return
	}

	optimizedInfo, err := os.Stat(tmpPath)
	if err != nil {
		os.Remove(tmpPath)
		return
	}

	if optimizedInfo.Size() < originalInfo.Size() {
		savedPct := float64(originalInfo.Size()-optimizedInfo.Size()) / float64(originalInfo.Size()) * 100
		log.Printf("[PDF] Optimizado: %s (%.1f%% reducción, %d → %d bytes)",
			filepath.Base(filePath), savedPct, originalInfo.Size(), optimizedInfo.Size())

		if err := os.Rename(tmpPath, filePath); err != nil {
			os.Remove(tmpPath)
			log.Printf("[PDF] Error al reemplazar archivo: %v", err)
		}
		return
	}

	os.Remove(tmpPath)
	log.Printf("[PDF] Sin reducción para %s, se conserva el original", filepath.Base(filePath))
}
