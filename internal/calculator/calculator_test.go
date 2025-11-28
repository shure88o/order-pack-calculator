package calculator

import (
	"reflect"
	"testing"
)

func TestCalculatePacks(t *testing.T) {
	tests := []struct {
		name      string
		packSizes []int
		orderQty  int
		wantPacks map[int]int
		wantTotal int
	}{
		{
			name:      "order 1 item",
			packSizes: []int{250, 500, 1000, 2000, 5000},
			orderQty:  1,
			wantPacks: map[int]int{250: 1},
			wantTotal: 250,
		},
		{
			name:      "exact pack size",
			packSizes: []int{250, 500, 1000, 2000, 5000},
			orderQty:  250,
			wantPacks: map[int]int{250: 1},
			wantTotal: 250,
		},
		{
			name:      "prefer fewer packs over more packs same items",
			packSizes: []int{250, 500, 1000, 2000, 5000},
			orderQty:  251,
			wantPacks: map[int]int{500: 1},
			wantTotal: 500,
		},
		{
			name:      "combination of packs",
			packSizes: []int{250, 500, 1000, 2000, 5000},
			orderQty:  501,
			wantPacks: map[int]int{500: 1, 250: 1},
			wantTotal: 750,
		},
		{
			name:      "large order",
			packSizes: []int{250, 500, 1000, 2000, 5000},
			orderQty:  12001,
			wantPacks: map[int]int{5000: 2, 2000: 1, 250: 1},
			wantTotal: 12250,
		},
		{
			name:      "order zero",
			packSizes: []int{250, 500, 1000, 2000, 5000},
			orderQty:  0,
			wantPacks: map[int]int{},
			wantTotal: 0,
		},
		{
			name:      "custom pack sizes",
			packSizes: []int{23, 31, 53},
			orderQty:  263,
			wantPacks: map[int]int{23: 2, 31: 7},
			wantTotal: 263,
		},
		{
			name:      "empty pack sizes list",
			packSizes: []int{},
			orderQty:  100,
			wantPacks: map[int]int{},
			wantTotal: 0,
		},
		{
			name:      "single pack size available",
			packSizes: []int{500},
			orderQty:  251,
			wantPacks: map[int]int{500: 1},
			wantTotal: 500,
		},
		{
			name:      "single pack size multiple needed",
			packSizes: []int{500},
			orderQty:  1200,
			wantPacks: map[int]int{500: 3},
			wantTotal: 1500,
		},
		{
			name:      "duplicate pack sizes in input",
			packSizes: []int{250, 500, 250, 1000},
			orderQty:  251,
			wantPacks: map[int]int{500: 1},
			wantTotal: 500,
		},
		{
			name:      "pack sizes without small values",
			packSizes: []int{1000, 5000},
			orderQty:  1,
			wantPacks: map[int]int{1000: 1},
			wantTotal: 1000,
		},
		{
			name:      "negative order quantity",
			packSizes: []int{250, 500},
			orderQty:  -5,
			wantPacks: map[int]int{},
			wantTotal: 0,
		},
		{
			name:      "order 500",
			packSizes: []int{250, 500, 1000, 2000, 5000},
			orderQty:  500,
			wantPacks: map[int]int{500: 1},
			wantTotal: 500,
		},
		{
			name:      "order 1000",
			packSizes: []int{250, 500, 1000, 2000, 5000},
			orderQty:  1000,
			wantPacks: map[int]int{1000: 1},
			wantTotal: 1000,
		},
		{
			name:      "order 1001",
			packSizes: []int{250, 500, 1000, 2000, 5000},
			orderQty:  1001,
			wantPacks: map[int]int{1000: 1, 250: 1},
			wantTotal: 1250,
		},
		{
			name:      "order 2500",
			packSizes: []int{250, 500, 1000, 2000, 5000},
			orderQty:  2500,
			wantPacks: map[int]int{2000: 1, 500: 1},
			wantTotal: 2500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPacks := CalculatePacks(tt.orderQty, tt.packSizes)

			// Calculate total items from result
			gotTotal := 0
			for size, qty := range gotPacks {
				gotTotal += size * qty
			}

			// Check if result matches expected
			if !reflect.DeepEqual(gotPacks, tt.wantPacks) {
				t.Errorf("CalculatePacks() packs = %v, want %v", gotPacks, tt.wantPacks)
			}

			if gotTotal != tt.wantTotal {
				t.Errorf("CalculatePacks() total = %d, want %d", gotTotal, tt.wantTotal)
			}
		})
	}
}

func TestRemoveDuplicates(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		want  []int
	}{
		{
			name:  "no duplicates",
			input: []int{1, 2, 3},
			want:  []int{1, 2, 3},
		},
		{
			name:  "with duplicates",
			input: []int{1, 2, 2, 3, 1},
			want:  []int{1, 2, 3},
		},
		{
			name:  "empty slice",
			input: []int{},
			want:  []int{},
		},
		{
			name:  "all same",
			input: []int{5, 5, 5},
			want:  []int{5},
		},
		{
			name:  "with zero and negative",
			input: []int{0, -1, 5, 0, 5},
			want:  []int{5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := removeDuplicates(tt.input)
			if !equalSlices(got, tt.want) {
				t.Errorf("removeDuplicates() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCopyMap(t *testing.T) {
	original := map[int]int{1: 2, 3: 4}
	copied := copyMap(original)

	// Modify original
	original[5] = 6

	// Check that copied wasn't affected
	if _, exists := copied[5]; exists {
		t.Error("copyMap() did not create a deep copy")
	}

	if len(copied) != 2 {
		t.Errorf("copyMap() length = %d, want 2", len(copied))
	}
}

// Helper function to compare slices ignoring order
func equalSlices(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}

	countA := make(map[int]int)
	countB := make(map[int]int)

	for _, v := range a {
		countA[v]++
	}
	for _, v := range b {
		countB[v]++
	}

	return reflect.DeepEqual(countA, countB)
}
