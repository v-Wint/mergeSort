package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"slices"
	"time"
)

/*
Generate random slice of integers
*/
func generateRandomSlice(size int, max int) []int {
	slice := make([]int, size)
	for i := 0; i < size; i++ {
		slice[i] = rand.Intn(max)
	}
	return slice
}

/*
Test sorting algorithm for correctness on slice
*/
func testSlice(a []int, f func([]int) []int) {
	fmt.Println("Input array:", a)
	r := f(a)
	fmt.Println("Result:", r)
	slices.Sort(a)
	fmt.Println("Correct:", slices.Equal(a, r), "\n")
}

/*
Test sorting algorithm
*/
func test(f func([]int) []int) {
	inputs := [][]int{
		make([]int, 0),
		make([]int, 10),
		make([]int, 10),
		make([]int, 10),
	}

	for i := 0; i < 10; i++ {
		inputs[1][i] = i
		inputs[2][i] = 10 - i
		inputs[3][i] = rand.Intn(100)
	}

	for i := 0; i < len(inputs); i++ {
		testSlice(inputs[i], f)
	}

	size := 1_000_000
	bigArray := generateRandomSlice(size, 100_000)
	fmt.Println("Big array", size, "elements")
	r := f(bigArray)
	slices.Sort(bigArray)
	fmt.Println("Correct:", slices.Equal(bigArray, r), "\n")
}

/*
Test sorting algorithm time on slices of different size
*/
func testTime(f func([]int) []int, size int, step int) {

	inputs := make([][]int, size)

	for i := 0; i < size; i++ {
		inputs[i] = make([]int, (i+1)*step)
	}

	for i := 0; i < size; i++ {
		for j := 0; j < len(inputs[i]); j++ {
			inputs[i][j] = rand.Intn(100000)
		}
	}

	results := make([]int, size)
	for i := 0; i < len(inputs); i++ {
		fmt.Println("Testing for", float64((i+1)*step)/1_000_000, "million")
		startTime := time.Now()
		f(inputs[i])
		r := time.Now().Sub(startTime)
		results[i] = int(r) / 1_000_000
		fmt.Println(r)
	}

	for i := range results {
		fmt.Printf("%6dM ", (i+1)*2)
	}
	fmt.Println()
	for _, v := range results {
		fmt.Printf("%6dms", v)
	}
}

/*
Get time of execution of a sorting algorithm f on a slice a
*/
func experiment(a []int, f func([]int) []int) int {
	for i := 0; i < 30; i++ {
		f(a)
	}
	var avg int
	for i := 0; i < 30; i++ {
		startTime := time.Now()
		f(a)
		avg += int(time.Now().Sub(startTime))
	}
	return avg / 30
}

/*
Experiment with different slice sizes and thresholds
*/
func experimentThreshold() {
	arraySizes := []int{
		500_000, 1_000_000, 2_000_000, 4_000_000, 8_000_000, 12_000_000,
	}

	thresholds := []int{
		8, 16, 32, 64, 128, 256, 512, 1024, 2048, 4096, 8192, 16384, 32768, 65536, 131072, 262144, 524288, 1048576, 2097152,
	}

	for _, arraySize := range arraySizes {
		fmt.Println("Testing sequential for arraySize", float64(arraySize)/1_000_000, "million")
		array := generateRandomSlice(arraySize, 100_000)
		t := experiment(array, mergeSort)
		fmt.Println("The result is", float64(t)/1_000_000, "ms\n")

		for _, threshold := range thresholds {
			fmt.Println("Testing concurrent for arraySize", float64(arraySize)/1_000_000,
				"million wit threshold", threshold)
			t := experiment(array, func(a []int) []int { return concurrentMergeSort(a, threshold) })
			fmt.Println("The result is:", float64(t)/1_000_000, "ms\n")
		}
	}
}

/*
Experiment with different number of cores allowed
*/
func experimentCores() {
	arraySizes := []int{
		500_000, 1_000_000, 2_000_000, 4_000_000, 8_000_000, 12_000_000, 16_000_000, 20_000_000, 32_000_000,
	}

	cores := []int{
		2, 4, 8,
	}

	for _, arraySize := range arraySizes {
		fmt.Println("Testing sequential for arraySize", float64(arraySize)/1_000_000, "million")
		array := generateRandomSlice(arraySize, 100_000)
		t := experiment(array, mergeSort)
		fmt.Println("The result is:", float64(t)/1_000_000, "ms\n")
		for _, core := range cores {
			fmt.Println("Testing concurrent for arraySize", float64(arraySize)/1_000_000,
				"million with", core, "cores")
			runtime.GOMAXPROCS(core)
			t := experiment(array, func(a []int) []int { return concurrentMergeSort(a, arraySize/core) })
			fmt.Println("The result is:", float64(t)/1_000_000, "ms\n")
		}
	}
}
