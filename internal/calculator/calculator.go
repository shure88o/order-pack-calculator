package calculator

import (
	"sort"
)

// Solution represents a pack combination with its metrics
type Solution struct {
	TotalItems int
	PackCount  int
	Packs      map[int]int // pack size -> quantity
}

// CalculatePacks calculates the optimal pack combination for an order
// Rules (in priority):
// 1. Only whole packs
// 2. Minimize excess items (total items - order quantity)
// 3. Minimize pack count (among solutions with same total items)
func CalculatePacks(orderQty int, packSizes []int) map[int]int {
	// Edge case: zero order returns empty result
	if orderQty <= 0 {
		return make(map[int]int)
	}

	// Edge case: no pack sizes available
	if len(packSizes) == 0 {
		return make(map[int]int)
	}

	// Remove duplicates and sort in descending order for optimization
	uniqueSizes := removeDuplicates(packSizes)
	sort.Sort(sort.Reverse(sort.IntSlice(uniqueSizes)))

	// Find the maximum pack size to determine upper bound
	maxPack := uniqueSizes[0]

	// Upper bound: we never need more excess than the largest pack
	upperBound := orderQty + maxPack

	// dp[i] stores the best solution to reach exactly i items
	dp := make([]*Solution, upperBound+1)

	// Base case: 0 items requires 0 packs
	dp[0] = &Solution{
		TotalItems: 0,
		PackCount:  0,
		Packs:      make(map[int]int),
	}

	// Fill DP table
	for amount := 1; amount <= upperBound; amount++ {
		for _, packSize := range uniqueSizes {
			if amount >= packSize && dp[amount-packSize] != nil {
				// Create candidate solution by adding one pack
				candidate := &Solution{
					TotalItems: amount,
					PackCount:  dp[amount-packSize].PackCount + 1,
					Packs:      copyMap(dp[amount-packSize].Packs),
				}
				candidate.Packs[packSize]++

				// Update dp[amount] if this candidate is better
				if dp[amount] == nil || candidate.PackCount < dp[amount].PackCount {
					dp[amount] = candidate
				}
			}
		}
	}

	// Find best solution: minimum items >= orderQty, then minimum packs
	var bestSolution *Solution
	for amount := orderQty; amount <= upperBound; amount++ {
		if dp[amount] != nil {
			if bestSolution == nil ||
				amount < bestSolution.TotalItems ||
				(amount == bestSolution.TotalItems && dp[amount].PackCount < bestSolution.PackCount) {
				bestSolution = dp[amount]
			}
		}
	}

	if bestSolution == nil {
		return make(map[int]int)
	}

	return bestSolution.Packs
}

// removeDuplicates removes duplicate pack sizes
func removeDuplicates(sizes []int) []int {
	seen := make(map[int]bool)
	result := []int{}

	for _, size := range sizes {
		if size > 0 && !seen[size] {
			seen[size] = true
			result = append(result, size)
		}
	}

	return result
}

// copyMap creates a deep copy of a map
func copyMap(m map[int]int) map[int]int {
	result := make(map[int]int)
	for k, v := range m {
		result[k] = v
	}
	return result
}
