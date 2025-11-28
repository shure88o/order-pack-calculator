package handler

import (
	"encoding/json"
	"net/http"
	"order-pack-calculator/internal/calculator"
	"order-pack-calculator/internal/model"
	"sort"
	"sync"
)

// Handler manages HTTP endpoints and pack configuration
type Handler struct {
	packSizes []int
	mu        sync.RWMutex
}

// NewHandler creates a new handler with initial pack sizes
func NewHandler(initialPackSizes []int) *Handler {
	return &Handler{
		packSizes: initialPackSizes,
	}
}

// GetPackSizes returns current pack sizes
func (h *Handler) GetPackSizes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	h.mu.RLock()
	sizes := make([]int, len(h.packSizes))
	copy(sizes, h.packSizes)
	h.mu.RUnlock()

	response := model.PackSizesResponse{
		PackSizes: sizes,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UpdatePackSizes updates pack sizes configuration
func (h *Handler) UpdatePackSizes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req model.PackSizesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.PackSizes) == 0 {
		sendError(w, "Pack sizes cannot be empty", http.StatusBadRequest)
		return
	}

	// Validate pack sizes (must be positive)
	for _, size := range req.PackSizes {
		if size <= 0 {
			sendError(w, "Pack sizes must be positive integers", http.StatusBadRequest)
			return
		}
	}

	h.mu.Lock()
	h.packSizes = req.PackSizes
	h.mu.Unlock()

	response := model.PackSizesResponse{
		PackSizes: req.PackSizes,
		Message:   "Pack sizes updated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CalculatePacks calculates optimal pack combination
func (h *Handler) CalculatePacks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req model.CalculateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.OrderQuantity < 0 {
		sendError(w, "Order quantity must be a non-negative integer", http.StatusBadRequest)
		return
	}

	h.mu.RLock()
	sizes := make([]int, len(h.packSizes))
	copy(sizes, h.packSizes)
	h.mu.RUnlock()

	// Calculate optimal packs
	packsMap := calculator.CalculatePacks(req.OrderQuantity, sizes)

	// Convert to response format
	packs := []model.PackBreakdown{}
	totalItems := 0
	totalPacks := 0

	// Sort by pack size descending for consistent output
	packSizes := make([]int, 0, len(packsMap))
	for size := range packsMap {
		packSizes = append(packSizes, size)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(packSizes)))

	for _, size := range packSizes {
		qty := packsMap[size]
		packs = append(packs, model.PackBreakdown{
			Size:     size,
			Quantity: qty,
		})
		totalItems += size * qty
		totalPacks += qty
	}

	response := model.CalculateResponse{
		OrderQuantity: req.OrderQuantity,
		Packs:         packs,
		TotalItems:    totalItems,
		TotalPacks:    totalPacks,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// sendError sends an error response
func sendError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(model.ErrorResponse{Error: message})
}
