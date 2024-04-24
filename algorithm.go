package main

import "sync"

func merge(left, right []int) []int {
	result := make([]int, 0, len(left)+len(right))

	for len(left) > 0 && len(right) > 0 {
		if left[0] < right[0] {
			result = append(result, left[0])
			left = left[1:]
		} else {
			result = append(result, right[0])
			right = right[1:]
		}
	}
	result = append(result, left...)
	result = append(result, right...)

	return result
}

func mergeSort(a []int) []int {
	if len(a) < 2 {
		return a
	}

	mid := len(a) / 2
	left := mergeSort(a[:mid])
	right := mergeSort(a[mid:])

	return merge(left, right)
}

func concurrentMergeSort(a []int, threshold int) []int {
	if len(a) <= threshold {
		return mergeSort(a)
	}

	mid := len(a) / 2
	leftChan := make(chan []int)
	rightChan := make(chan []int)

	go func() { leftChan <- concurrentMergeSort(a[:mid], threshold) }()
	go func() { rightChan <- concurrentMergeSort(a[mid:], threshold) }()

	return merge(<-leftChan, <-rightChan)
}

func waitGroupConcurrentMergeSort(a []int, threshold int) []int {
	if len(a) <= threshold {
		return mergeSort(a)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	mid := len(a) / 2

	var left []int
	var right []int

	go func() { defer wg.Done(); left = waitGroupConcurrentMergeSort(a[:mid], threshold) }()
	go func() { defer wg.Done(); right = waitGroupConcurrentMergeSort(a[mid:], threshold) }()

	wg.Wait()

	return merge(left, right)
}
