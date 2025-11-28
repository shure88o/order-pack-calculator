package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"order-pack-calculator/internal/model"
	"testing"
)

func TestGetPackSizes(t *testing.T) {
	handler := NewHandler([]int{250, 500, 1000})

	req := httptest.NewRequest(http.MethodGet, "/api/packs", nil)
	w := httptest.NewRecorder()

	handler.GetPackSizes(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response model.PackSizesResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(response.PackSizes) != 3 {
		t.Errorf("Expected 3 pack sizes, got %d", len(response.PackSizes))
	}
}

func TestUpdatePackSizes(t *testing.T) {
	handler := NewHandler([]int{250, 500})

	reqBody := model.PackSizesRequest{
		PackSizes: []int{100, 200, 300},
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPut, "/api/packs", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.UpdatePackSizes(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Verify pack sizes were updated
	handler.mu.RLock()
	if len(handler.packSizes) != 3 {
		t.Errorf("Expected 3 pack sizes, got %d", len(handler.packSizes))
	}
	handler.mu.RUnlock()
}

func TestUpdatePackSizesInvalid(t *testing.T) {
	handler := NewHandler([]int{250, 500})

	tests := []struct {
		name    string
		request model.PackSizesRequest
	}{
		{
			name:    "empty pack sizes",
			request: model.PackSizesRequest{PackSizes: []int{}},
		},
		{
			name:    "negative pack size",
			request: model.PackSizesRequest{PackSizes: []int{250, -100}},
		},
		{
			name:    "zero pack size",
			request: model.PackSizesRequest{PackSizes: []int{0, 100}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.request)
			req := httptest.NewRequest(http.MethodPut, "/api/packs", bytes.NewReader(body))
			w := httptest.NewRecorder()

			handler.UpdatePackSizes(w, req)

			if w.Code != http.StatusBadRequest {
				t.Errorf("Expected status 400, got %d", w.Code)
			}
		})
	}
}

func TestCalculatePacks(t *testing.T) {
	handler := NewHandler([]int{250, 500, 1000, 2000, 5000})

	reqBody := model.CalculateRequest{
		OrderQuantity: 251,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/calculate", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.CalculatePacks(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response model.CalculateResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.OrderQuantity != 251 {
		t.Errorf("Expected order quantity 251, got %d", response.OrderQuantity)
	}

	if response.TotalItems != 500 {
		t.Errorf("Expected total items 500, got %d", response.TotalItems)
	}

	if response.TotalPacks != 1 {
		t.Errorf("Expected 1 pack, got %d", response.TotalPacks)
	}
}

func TestCalculatePacksInvalid(t *testing.T) {
	handler := NewHandler([]int{250, 500})

	reqBody := model.CalculateRequest{
		OrderQuantity: -5,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/calculate", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.CalculatePacks(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestPackSizesPersistence(t *testing.T) {
	handler := NewHandler([]int{250, 500, 1000})

	// Update pack sizes
	updateReq := model.PackSizesRequest{
		PackSizes: []int{23, 31, 53},
	}
	updateBody, _ := json.Marshal(updateReq)
	updateHTTP := httptest.NewRequest(http.MethodPut, "/api/packs", bytes.NewReader(updateBody))
	updateW := httptest.NewRecorder()
	handler.UpdatePackSizes(updateW, updateHTTP)

	// Calculate with new pack sizes
	calcReq := model.CalculateRequest{
		OrderQuantity: 263,
	}
	calcBody, _ := json.Marshal(calcReq)
	calcHTTP := httptest.NewRequest(http.MethodPost, "/api/calculate", bytes.NewReader(calcBody))
	calcW := httptest.NewRecorder()
	handler.CalculatePacks(calcW, calcHTTP)

	var calcResponse model.CalculateResponse
	if err := json.NewDecoder(calcW.Body).Decode(&calcResponse); err != nil {
		t.Fatalf("Failed to decode calculate response: %v", err)
	}

	// Verify calculation uses new pack sizes (should get exact match at 263)
	if calcResponse.TotalItems != 263 {
		t.Errorf("Expected calculation with new pack sizes to give 263 items, got %d", calcResponse.TotalItems)
	}
}
