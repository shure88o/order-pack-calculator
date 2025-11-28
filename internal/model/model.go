package model

// PackSizesRequest represents a request to update pack sizes
type PackSizesRequest struct {
	PackSizes []int `json:"pack_sizes"`
}

// PackSizesResponse represents pack sizes data
type PackSizesResponse struct {
	PackSizes []int  `json:"pack_sizes"`
	Message   string `json:"message,omitempty"`
}

// CalculateRequest represents a request to calculate optimal packs
type CalculateRequest struct {
	OrderQuantity int `json:"order_quantity"`
}

// PackBreakdown represents a single pack size and its quantity
type PackBreakdown struct {
	Size     int `json:"size"`
	Quantity int `json:"quantity"`
}

// CalculateResponse represents the result of pack calculation
type CalculateResponse struct {
	OrderQuantity int             `json:"order_quantity"`
	Packs         []PackBreakdown `json:"packs"`
	TotalItems    int             `json:"total_items"`
	TotalPacks    int             `json:"total_packs"`
}

// ErrorResponse represents an error message
type ErrorResponse struct {
	Error string `json:"error"`
}
