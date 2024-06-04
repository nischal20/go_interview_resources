package main

import (
	"fmt"
)

// Function to find the length of the Longest Common Subsequence
func lcs(X, Y string) int {
	m := len(X)
	n := len(Y)

	// Create a 2D slice to store lengths of LCS
	L := make([][]int, m+1)
	for i := range L {
		L[i] = make([]int, n+1)
	}

	// Build the L[m+1][n+1] in bottom-up fashion
	for i := 0; i <= m; i++ {
		for j := 0; j <= n; j++ {
			if i == 0 || j == 0 {
				L[i][j] = 0
			} else if X[i-1] == Y[j-1] {
				L[i][j] = L[i-1][j-1] + 1
			} else {
				L[i][j] = max(L[i-1][j], L[i][j-1])
			}
		}
	}

	return L[m][n]
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	X := "AGGTAB"
	Y := "GXTXAYB"

	fmt.Printf("Length of LCS is %d\n", lcs(X, Y))
}
