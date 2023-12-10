package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

type RequestBody struct {
	ToSort [][]int `json:"to_sort"`
}

type ResponseBody struct {
	SortedArrays [][]int `json:"sorted_arrays"`
	TimeNS       float64 `json:"time_ns"`
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/process-single", handleSort).Methods("POST")
	r.HandleFunc("/process-concurrent", handleConcurrentSort).Methods("POST")

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

func mergeSort(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}

	mid := len(arr) / 2
	left := mergeSort(arr[:mid])
	right := mergeSort(arr[mid:])

	return merge(left, right)
}

func merge(left, right []int) []int {
	result := make([]int, 0, len(left)+len(right))

	for len(left) > 0 || len(right) > 0 {
		if len(left) == 0 {
			return append(result, right...)
		}
		if len(right) == 0 {
			return append(result, left...)
		}

		if left[0] <= right[0] {
			result = append(result, left[0])
			left = left[1:]
		} else {
			result = append(result, right[0])
			right = right[1:]
		}
	}

	return result
}

func handleSort(w http.ResponseWriter, r *http.Request) {
	var requestBody RequestBody
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	startTime := time.Now()
	fmt.Printf(startTime.String())
	sortedArrays := make([][]int, len(requestBody.ToSort))
	for i, arr := range requestBody.ToSort {
		sortedArrays[i] = mergeSort(arr)
	}

	elapsedTime := time.Since(startTime).Seconds()

	responseBody := ResponseBody{
		SortedArrays: sortedArrays,
		TimeNS:       elapsedTime * 1e9,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(responseBody)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
}

func handleConcurrentSort(w http.ResponseWriter, r *http.Request) {
	var requestBody RequestBody
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	startTime := time.Now()
	fmt.Printf(startTime.String())
	var wg sync.WaitGroup
	wg.Add(len(requestBody.ToSort))
	resultCh := make(chan struct {
		index int
		data  []int
	}, len(requestBody.ToSort))

	for i, arr := range requestBody.ToSort {
		go func(index int, arr []int) {
			defer wg.Done()
			resultCh <- struct {
				index int
				data  []int
			}{index, mergeSort(arr)}
		}(i, arr)
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	sortedArrays := make([][]int, len(requestBody.ToSort))
	for result := range resultCh {
		sortedArrays[result.index] = result.data
	}

	elapsedTime := time.Since(startTime).Seconds()
	responseBody := ResponseBody{
		SortedArrays: sortedArrays,
		TimeNS:       elapsedTime * 1e9,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(responseBody)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
}
