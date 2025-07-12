package main

import "sync"

const threshold = 1_000_000 // Don't spawn goroutines for tiny chunks
func gomergeSort(arr []int, left int, right int, wg *sync.WaitGroup) {

	if left >= right {
		return
	}
	mid := (left + right) / 2
	if right-left > threshold {

		var wglocal sync.WaitGroup
		wglocal.Add(2)
		wg.Add(1)

		defer wg.Done()
		go func() {
			defer wglocal.Done()
			gomergeSort(arr, left, mid, wg)
		}()

		go func() {
			defer wglocal.Done()
			gomergeSort(arr, mid+1, right, wg)

		}()
		wglocal.Wait()
	} else {
		gomergeSort(arr, left, mid, wg)
		gomergeSort(arr, mid+1, right, wg)
	}
	merge(arr, left, mid, right)
}
