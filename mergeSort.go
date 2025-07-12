package main

func mergeSort(arr []int, left int, right int) {
	if left >= right {
		return
	}
	mid := (left + right) / 2

	mergeSort(arr, left, mid)
	mergeSort(arr, mid+1, right)

	merge(arr, left, mid, right)
}
func merge(arr []int, left int, mid int, right int) {
	temp := make([]int, 0, right-left+1)
	i := left
	j := mid + 1
	for i <= mid && j <= right {
		if arr[i] <= arr[j] {
			temp = append(temp, arr[i])
			i++
		} else {
			temp = append(temp, arr[j])
			j++
		}
	}
	for i <= mid {
		temp = append(temp, arr[i])
		i++
	}
	for j <= right {
		temp = append(temp, arr[j])
		j++
	}
	for k := left; k <= right; k++ {
		arr[k] = temp[k-left]

	}
}
