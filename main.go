// package main

// import (
// 	"fmt"
// 	"math/rand"

// 	_ "net/http/pprof"
// 	"runtime"

// 	"sync"
// 	"time"
// )

// var SIZES = []int{1e7}

// func main() {
// 	// f, _ := os.Create("cpu.prof")
// 	// pprof.StartCPUProfile(f)

// 	fmt.Println("Write the size of array")
// 	for _, size := range SIZES {

// 		arr := make([]int, 0, size)
// 		for i := 0; i < size; i++ {
// 			arr = append(arr, rand.Intn(size))
// 		}
// 		fmt.Println("Size of the array after adding elements")
// 		// Prepare copies
// 		arr1 := make([]int, size)
// 		arr2 := make([]int, size)
// 		copy(arr1, arr)
// 		copy(arr2, arr)

// 		// fmt.Println("Initial Stats Before Any Sort")
// 		// printMemStats("INITIAL")

// 		// // --- Standard Merge Sort ---
// 		// fmt.Println("\n\n--- Running Standard Merge Sort ---")
// 		// runtime.GC()
// 		printMemStats("Before Standard Merge Sort")
// 		start := time.Now()
// 		mergeSort(arr1, 0, len(arr1)-1)
// 		duration := time.Since(start)
// 		fmt.Printf("Standard Merge Sort took: %s\n", duration)

// 		// --- Concurrent Merge Sort ---
// 		start = time.Now()
// 		var wg sync.WaitGroup
// 		// wg.Add(1)
// 		gomergeSort(arr2, 0, len(arr2)-1, &wg)
// 		wg.Wait()
// 		duration = time.Since(start)
// 		fmt.Printf("Concurrent Merge Sort took: %s\n", duration)
// 	}
// 	// pprof.StopCPUProfile()
// }

//	func printMemStats(label string) {
//		var m runtime.MemStats
//		runtime.ReadMemStats(&m)
//		fmt.Println("========== " + label + " ==========")
//		fmt.Printf("Alloc = %.2f MB\n", float64(m.Alloc)/1024/1024)
//		fmt.Printf("TotalAlloc = %.2f MB\n", float64(m.TotalAlloc)/1024/1024)
//		fmt.Printf("Sys = %.2f MB\n", float64(m.Sys)/1024/1024)
//		fmt.Printf("HeapAlloc = %.2f MB\n", float64(m.HeapAlloc)/1024/1024)
//		fmt.Printf("HeapSys = %.2f MB\n", float64(m.HeapSys)/1024/1024)
//		fmt.Printf("NumGC = %v\n", m.NumGC)
//		fmt.Printf("NumGoroutine = %v\n", runtime.NumGoroutine())
//		fmt.Printf("GOMAXPROCS = %v\n", runtime.GOMAXPROCS(0))
//		fmt.Printf("NumCPU = %v\n", runtime.NumCPU())
//		fmt.Println("===================================")
//	}
package main

import (
	"fmt"
	"math"

	"math/rand"
	"runtime"
	"sync"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

var SIZES = []int{1e4, 1e5, 1e6, 1e7}

type Stat struct {
	Size      int
	TimeStd   float64
	TimeCon   float64
	NumGC     uint32
	HeapAlloc float64
	Goroutine int
}

func main() {
	var stats []Stat

	for _, size := range SIZES {
		fmt.Printf("Benchmarking size: %d\n", size)

		arr := make([]int, 0, size)
		for i := 0; i < size; i++ {
			arr = append(arr, rand.Intn(size))
		}
		arr1 := make([]int, size)
		arr2 := make([]int, size)
		copy(arr1, arr)
		copy(arr2, arr)

		// --- Standard Merge Sort ---
		runtime.GC()
		start := time.Now()
		mergeSort(arr1, 0, len(arr1)-1)
		durationStd := time.Since(start)

		// --- Concurrent Merge Sort ---
		runtime.GC()
		start = time.Now()
		var wg sync.WaitGroup
		gomergeSort(arr2, 0, len(arr2)-1, &wg)
		wg.Wait()
		durationCon := time.Since(start)

		// --- Collect Memory Stats ---
		var mem runtime.MemStats
		runtime.ReadMemStats(&mem)

		stats = append(stats, Stat{
			Size:      size,
			TimeStd:   float64(durationStd.Milliseconds()),
			TimeCon:   float64(durationCon.Milliseconds()),
			NumGC:     mem.NumGC,
			HeapAlloc: float64(mem.HeapAlloc) / 1024.0 / 1024.0, // MB
			Goroutine: runtime.NumGoroutine(),
		})
	}

	plotTiming(stats)
	plotMemory(stats)
	fmt.Println("Plots generated: plot.png, mem_plot.png")
}

// ---------------- PLOTTING ----------------

func plotTiming(stats []Stat) {
	p := plot.New()
	p.Title.Text = "Merge Sort Time Comparison"
	p.X.Label.Text = "Input Size"
	p.Y.Label.Text = "Time (ms)"
	p.Y.Scale = plot.LogScale{}
	p.Y.Tick.Marker = plot.LogTicks{}

	stdLine := make(plotter.XYs, len(stats))
	conLine := make(plotter.XYs, len(stats))

	for i, s := range stats {
		stdLine[i].X = float64(s.Size)
		stdLine[i].Y = math.Max(s.TimeStd, 0.1)
		conLine[i].X = float64(s.Size)
		conLine[i].Y = math.Max(s.TimeCon, 0.1)
	}

	err := plotutil.AddLinePoints(p,
		"Standard Merge Sort", stdLine,
		"Concurrent Merge Sort", conLine,
	)
	if err != nil {
		panic(err)
	}

	if err := p.Save(6*vg.Inch, 4*vg.Inch, "plot.png"); err != nil {
		panic(err)
	}
}
func plotMemory(stats []Stat) {
	p := plot.New()
	p.Title.Text = "Memory & Runtime Stats"
	p.X.Label.Text = "Input Size"
	p.Y.Label.Text = "Metric Value"

	gcLine := make(plotter.XYs, len(stats))
	heapLine := make(plotter.XYs, len(stats))
	gorLine := make(plotter.XYs, len(stats))

	for i, s := range stats {
		x := float64(s.Size)
		gcLine[i].X = x
		gcLine[i].Y = float64(s.NumGC)
		heapLine[i].X = x
		heapLine[i].Y = s.HeapAlloc // MB
		gorLine[i].X = x
		gorLine[i].Y = float64(s.Goroutine)
	}

	err := plotutil.AddLinePoints(p,
		"GC Count (NumGC)", gcLine,
		"Heap Alloc (MB)", heapLine,
		"Goroutines", gorLine,
	)
	if err != nil {
		panic(err)
	}

	// Format X-axis ticks
	p.X.Tick.Marker = plot.ConstantTicks([]plot.Tick{
		{Value: 1e4, Label: "10K"},
		{Value: 1e5, Label: "100K"},
		{Value: 1e6, Label: "1M"},
		{Value: 1e7, Label: "10M"},
	})

	if err := p.Save(6*vg.Inch, 4*vg.Inch, "mem_plot.png"); err != nil {
		panic(err)
	}
}
