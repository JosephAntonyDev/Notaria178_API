package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

type LocalFileStorage struct {
	BasePath string
}

func NewLocalFileStorage() *LocalFileStorage {
	var base string
	if runtime.GOOS == "windows" {
		base = `C:\Users\MSI\Documents\Notaria178\uploads`
	} else {
		base = "/var/uploads/notaria178"
	}
	return &LocalFileStorage{BasePath: base}
}

// SaveFile guarda el archivo físico en disco y devuelve la ruta absoluta resultante.
// Estructura: BasePath/branches/{branchID}/works/{workID}/{timestamp}_{filename}
func (s *LocalFileStorage) SaveFile(file *multipart.FileHeader, branchID, workID string) (string, error) {
	dir := filepath.Join(s.BasePath, "branches", branchID, "works", workID)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("error al crear directorio de almacenamiento: %w", err)
	}

	timestamp := time.Now().Format("20060102150405")
	safeName := filepath.Base(file.Filename)
	fileName := fmt.Sprintf("%s_%s", timestamp, safeName)
	fullPath := filepath.Join(dir, fileName)

	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("error al abrir el archivo subido: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("error al crear el archivo en disco: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("error al copiar el archivo: %w", err)
	}

	// Cerrar el archivo destino antes de comprimir
	dst.Close()

	// Compresión lossless para PDFs
	CompressPDF(fullPath)

	return fullPath, nil
}
