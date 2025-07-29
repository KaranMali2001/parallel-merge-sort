package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime/pprof"

	_ "net/http/pprof"
	"runtime"

	"sync"
	"time"
)

var SIZES = []int{1e7, 1e8}

func main() {
	f, _ := os.Create("cpu.prof")
	pprof.StartCPUProfile(f)

	fmt.Println("Write the size of array")
	for _, size := range SIZES {

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

		// fmt.Println("Initial Stats Before Any Sort")
		// printMemStats("INITIAL")

		// // --- Standard Merge Sort ---
		// fmt.Println("\n\n--- Running Standard Merge Sort ---")
		// runtime.GC()
		printMemStats("Before Standard Merge Sort")
		start := time.Now()
		mergeSort(arr1, 0, len(arr1)-1)
		duration := time.Since(start)
		fmt.Printf("Standard Merge Sort took: %s\n", duration)

		// --- Concurrent Merge Sort ---
		start = time.Now()
		var wg sync.WaitGroup
		// wg.Add(1)
		gomergeSort(arr2, 0, len(arr2)-1, &wg)
		wg.Wait()
		duration = time.Since(start)
		fmt.Printf("Concurrent Merge Sort took: %s\n", duration)
	}
	pprof.StopCPUProfile()
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
