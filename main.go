package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strconv"

	"sync"
	"time"
)

func main() {
	fmt.Println("Write the size of array")
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./app <array_size>")
		os.Exit(1)
	}
	size, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic("enter Number in dockerfile")
	}
	arr := make([]int, 0, size)
	for i := 0; i < size; i++ {
		arr = append(arr, rand.Intn(size))
	}
	fmt.Println("Size of the array after adding elements")
	// Prepare copies
	arr1 := make([]int, size)
	arr2 := make([]int, size)
	copy(arr1, arr)
	copy(arr2, arr)

	fmt.Println("Initial Stats Before Any Sort")
	printMemStats("INITIAL")

	// --- Standard Merge Sort ---
	fmt.Println("\n\n--- Running Standard Merge Sort ---")
	runtime.GC()
	printMemStats("Before Standard Merge Sort")
	start := time.Now()
	mergeSort(arr1, 0, len(arr1)-1)
	duration := time.Since(start)
	printMemStats("After Standard Merge Sort")
	fmt.Printf("Standard Merge Sort took: %s\n", duration)
	fmt.Printf("Sorted: %v ... %v\n", arr1[:5], arr1[len(arr1)-5:])

	// --- Concurrent Merge Sort ---
	fmt.Println("\n\n--- Running Concurrent Merge Sort ---")
	runtime.GC()
	printMemStats("Before Concurrent Merge Sort")
	start = time.Now()
	var wg sync.WaitGroup
	// wg.Add(1)
	gomergeSort(arr2, 0, len(arr2)-1, &wg)
	wg.Wait()
	duration = time.Since(start)
	printMemStats("After Concurrent Merge Sort")
	fmt.Printf("Concurrent Merge Sort took: %s\n", duration)
	fmt.Printf("Sorted: %v ... %v\n", arr2[:5], arr2[len(arr2)-5:])
}

func printMemStats(label string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Println("========== " + label + " ==========")
	fmt.Printf("Alloc = %.2f MB\n", float64(m.Alloc)/1024/1024)
	fmt.Printf("TotalAlloc = %.2f MB\n", float64(m.TotalAlloc)/1024/1024)
	fmt.Printf("Sys = %.2f MB\n", float64(m.Sys)/1024/1024)
	fmt.Printf("HeapAlloc = %.2f MB\n", float64(m.HeapAlloc)/1024/1024)
	fmt.Printf("HeapSys = %.2f MB\n", float64(m.HeapSys)/1024/1024)
	fmt.Printf("NumGC = %v\n", m.NumGC)
	fmt.Printf("NumGoroutine = %v\n", runtime.NumGoroutine())
	fmt.Printf("GOMAXPROCS = %v\n", runtime.GOMAXPROCS(0))
	fmt.Printf("NumCPU = %v\n", runtime.NumCPU())
	fmt.Println("===================================")
}
