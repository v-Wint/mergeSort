package main

func main() {
	test(mergeSort)
	testTime(mergeSort, 20, 2_000_000)

	test(func(a []int) []int {
		return concurrentMergeSort(a, 2048)
	})
	testTime(func(a []int) []int {
		return concurrentMergeSort(a, 250_000)
	}, 20, 2_000_000)
	experimentThreshold()
	experimentCores()
}
