package main

import (
	"log"
	"net/http"
	"order-pack-calculator/internal/handler"
	"os"
	"strconv"
	"strings"
)

func main() {
	port := getEnv("PORT", "8080")
	packSizes := parsePackSizes(getEnv("PACK_SIZES", "250,500,1000,2000,5000"))

	h := handler.NewHandler(packSizes)

	// API routes
	http.HandleFunc("/api/packs", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w)
		if r.Method == http.MethodOptions {
			return
		}
		if r.Method == http.MethodGet {
			h.GetPackSizes(w, r)
		} else if r.Method == http.MethodPut {
			h.UpdatePackSizes(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/calculate", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w)
		if r.Method == http.MethodOptions {
			return
		}
		h.CalculatePacks(w, r)
	})

	// Serve static files and frontend
	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/", fs)

	addr := ":" + port
	log.Printf("Server starting on http://0.0.0.0%s", addr)
	log.Printf("Initial pack sizes: %v", packSizes)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// parsePackSizes parses comma-separated pack sizes
func parsePackSizes(s string) []int {
	parts := strings.Split(s, ",")
	sizes := make([]int, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if num, err := strconv.Atoi(part); err == nil && num > 0 {
			sizes = append(sizes, num)
		}
	}

	// Default if parsing failed
	if len(sizes) == 0 {
		return []int{250, 500, 1000, 2000, 5000}
	}

	return sizes
}

// enableCORS adds CORS headers
func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}
