package services

import "math"

type dpEntry struct {
	count int
	prev  int
	pack  int
}

// GetPacks calculates the minimal number of packs to achieve a certain amount
// and returns the count of each pack size.
func GetPacks(N int, packSizes []int) map[int]int {
	if len(packSizes) == 0 {
		return nil
	}

	maxCheck := N + packSizes[0]
	dp := initializeDP(maxCheck)

	// Fill DP table with minimal pack counts
	fillDP(dp, N, maxCheck, packSizes)

	// Find the best sum with the minimum leftover and pack count
	bestSum := findBestSum(dp, N, maxCheck)

	if bestSum == -1 {
		return nil
	}

	// Reconstruct the solution from the DP table
	return reconstructSolution(dp, bestSum)
}

// initializeDP initializes the dp array to store pack counts, with a default of -1 for all values.
func initializeDP(maxCheck int) []dpEntry {
	dp := make([]dpEntry, maxCheck+1)
	for i := range dp {
		dp[i].count = -1
	}
	dp[0].count = 0
	return dp
}

// fillDP populates the dp table with the minimum number of packs for each sum
func fillDP(dp []dpEntry, N, maxCheck int, packSizes []int) {
	for x := 0; x <= maxCheck; x++ {
		if dp[x].count == -1 {
			continue
		}
		for _, p := range packSizes {
			next := x + p
			if next <= maxCheck && (dp[next].count == -1 || dp[next].count > dp[x].count+1) {
				dp[next].count = dp[x].count + 1
				dp[next].prev = x
				dp[next].pack = p
			}
		}
	}
}

// findBestSum identifies the best sum with minimal leftover and minimal pack count
func findBestSum(dp []dpEntry, N, maxCheck int) int {
	minLeftover := math.MaxInt
	minPacks := math.MaxInt
	bestSum := -1

	for x := N; x <= maxCheck; x++ {
		if dp[x].count != -1 {
			leftover := x - N
			if leftover < minLeftover || (leftover == minLeftover && dp[x].count < minPacks) {
				minLeftover = leftover
				minPacks = dp[x].count
				bestSum = x
			}
		}
	}

	return bestSum
}

// reconstructSolution constructs the pack count map from the dp table
func reconstructSolution(dp []dpEntry, bestSum int) map[int]int {
	packCount := make(map[int]int)
	current := bestSum
	for current > 0 {
		entry := dp[current]
		packCount[entry.pack]++
		current = entry.prev
	}
	return packCount
}
