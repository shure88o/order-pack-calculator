package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// Storage provides persistence for pack sizes configuration
type Storage struct {
	filename string
	mu       sync.RWMutex
}

// NewStorage creates a new storage instance
// filename: path to JSON file for storing pack sizes
func NewStorage(filename string) *Storage {
	return &Storage{
		filename: filename,
	}
}

// LoadPackSizes loads pack sizes from file
// Returns the pack sizes and any error encountered
func (s *Storage) LoadPackSizes() ([]int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Check if file exists
	if _, err := os.Stat(s.filename); os.IsNotExist(err) {
		// File doesn't exist, return empty slice (will use defaults)
		return nil, nil
	}

	// Read file
	data, err := os.ReadFile(s.filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read storage file: %w", err)
	}

	// Parse JSON
	var packSizes []int
	if len(data) == 0 {
		// Empty file, return nil
		return nil, nil
	}

	if err := json.Unmarshal(data, &packSizes); err != nil {
		return nil, fmt.Errorf("failed to parse storage file: %w", err)
	}

	return packSizes, nil
}

// SavePackSizes saves pack sizes to file
// Returns any error encountered during save
func (s *Storage) SavePackSizes(packSizes []int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Marshal to JSON
	data, err := json.MarshalIndent(packSizes, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal pack sizes: %w", err)
	}

	// Write to file
	if err := os.WriteFile(s.filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write storage file: %w", err)
	}

	return nil
}

// FileExists checks if the storage file exists
func (s *Storage) FileExists() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, err := os.Stat(s.filename)
	return !os.IsNotExist(err)
}
