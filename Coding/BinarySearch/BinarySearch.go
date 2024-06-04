package main

import "fmt"

func lowerBound(arr []int, target int) int {

	left := 0
	right := len(arr)

	for left < right {
		mid := left + (right-left)/2
		if arr[mid] < target {
			left = mid + 1
		} else {
			right = mid
		}
	}
	return right
}

func upperBound(arr []int, target int) int {

	left := 0
	right := len(arr)

	for left < right {
		mid := left + (right-left)/2
		if arr[mid] <= target {
			left = mid + 1
		} else {
			right = mid
		}
	}
	return right
}

func BinarySearch(arr []int, target int) int {

	left := 0
	right := len(arr)

	for left < right {
		mid := left + (right-left)/2
		if arr[mid] == target {
			return mid
		} else if arr[mid] < target {
			left = mid + 1
		} else {
			right = mid + 1
		}
	}
	return -1
}

func main() {
	arr := []int{1, 2, 3, 5, 5, 6, 6, 8, 8, 9, 9, 10}
	target := 5
	index := lowerBound(arr, target)
	fmt.Printf("The lower bound of %d in the array is at index %d\n", target, index)
	index = upperBound(arr, target)
	fmt.Printf("The upperBound bound of %d in the array is at index %d\n", target, index)
	index = BinarySearch(arr, target)
	fmt.Printf("The BinarySearch  of %d in the array is at index %d\n", target, index)

}
