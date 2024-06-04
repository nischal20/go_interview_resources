package main

import (
	"fmt"
)

// heapify function to maintain the min-heap property
func heapify(array []int, n int, i int) {
	smallest := i
	left := 2*i + 1
	right := 2*i + 2

	// If left child is smaller than root
	if left < n && array[left] < array[smallest] {
		smallest = left
	}

	// If right child is smaller than the smallest so far
	if right < n && array[right] < array[smallest] {
		smallest = right
	}

	// If smallest is not root
	if smallest != i {
		swap(array, i, smallest)
		// Recursively heapify the affected sub-tree
		heapify(array, n, smallest)
	}
}

// Utility function to swap two elements in the array
func swap(array []int, i, j int) {
	array[i], array[j] = array[j], array[i]
}

// Function to build a min-heap from an array
func buildMinHeap(array []int) {
	n := len(array)
	// Start from the last non-leaf node and go up to the root
	for i := n/2 - 1; i >= 0; i-- {
		heapify(array, n, i)
	}
}

func main() {
	array := []int{4, 10, 3, 5, 1}
	fmt.Println("Original array:", array)

	buildMinHeap(array)
	fmt.Println("Min-heap array:", array)
}
