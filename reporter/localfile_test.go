package reporter

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/future-architect/vuls/models"
)

func TestLocalFileWriter_PathTraversal(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "vuls-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	w := LocalFileWriter{
		CurrentDir: tempDir,
		FormatJSON: true,
		FormatCsv:  true,
	}

	// Normal result should succeed
	validResult := models.ScanResult{
		ServerName: "valid-server",
	}
	if err := w.Write(validResult); err != nil {
		t.Fatalf("expected valid write to succeed, got: %v", err)
	}

	// Result with path traversal in ServerName should fail
	traversalResult := models.ScanResult{
		ServerName: "../../evil",
	}
	if err := w.Write(traversalResult); err == nil {
		t.Fatalf("expected path traversal in ServerName to fail, but got nil")
	}

	// formatCsvList with path traversal should fail
	if err := formatCsvList(validResult, filepath.Join(tempDir, "../../../evil.csv")); err == nil {
		t.Fatalf("expected formatCsvList with traversal path to fail, but got nil")
	}
}
