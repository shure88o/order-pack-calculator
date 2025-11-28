package storage

import (
	"os"
	"testing"
)

func TestStorage(t *testing.T) {
	// Use temporary file for testing
	tmpFile := "test_pack_sizes.json"
	defer os.Remove(tmpFile) // Clean up

	storage := NewStorage(tmpFile)

	// Test saving pack sizes
	packSizes := []int{250, 500, 1000, 2000, 5000}
	if err := storage.SavePackSizes(packSizes); err != nil {
		t.Fatalf("Failed to save pack sizes: %v", err)
	}

	// Verify file exists
	if !storage.FileExists() {
		t.Error("Storage file should exist after saving")
	}

	// Test loading pack sizes
	loaded, err := storage.LoadPackSizes()
	if err != nil {
		t.Fatalf("Failed to load pack sizes: %v", err)
	}

	if len(loaded) != len(packSizes) {
		t.Errorf("Loaded pack sizes length mismatch: got %d, want %d", len(loaded), len(packSizes))
	}

	for i, size := range packSizes {
		if loaded[i] != size {
			t.Errorf("Pack size mismatch at index %d: got %d, want %d", i, loaded[i], size)
		}
	}
}

func TestStorageNonExistentFile(t *testing.T) {
	storage := NewStorage("non_existent_file.json")

	// Loading from non-existent file should return nil, nil
	packSizes, err := storage.LoadPackSizes()
	if err != nil {
		t.Errorf("Loading from non-existent file should not error, got: %v", err)
	}
	if packSizes != nil {
		t.Errorf("Loading from non-existent file should return nil, got: %v", packSizes)
	}
}

func TestStorageEmptyFile(t *testing.T) {
	tmpFile := "test_empty.json"
	defer os.Remove(tmpFile)

	// Create empty file
	if err := os.WriteFile(tmpFile, []byte(""), 0644); err != nil {
		t.Fatalf("Failed to create empty file: %v", err)
	}

	storage := NewStorage(tmpFile)
	packSizes, err := storage.LoadPackSizes()
	if err != nil {
		t.Errorf("Loading from empty file should not error, got: %v", err)
	}
	if packSizes != nil {
		t.Errorf("Loading from empty file should return nil, got: %v", packSizes)
	}
}
