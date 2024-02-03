package websocket

import (
	"encoding/json"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
)

func bubbleSortWithProgress(conn *websocket.Conn, data []int, getStepDelay func() time.Duration) {
	n := len(data)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if data[j] > data[j+1] {
				data[j], data[j+1] = data[j+1], data[j]
				sendUpdate(conn, data, j+1, getStepDelay)
			}
		}
	}
}

func quicksortWithProgress(conn *websocket.Conn, data []int, start, end int, getStepDelay func() time.Duration) {
	if start < end {
		pivot := data[end]
		i := start - 1

		for j := start; j < end; j++ {
			if data[j] <= pivot {
				i++
				data[i], data[j] = data[j], data[i]
				sendUpdate(conn, data, j, getStepDelay)
			}
		}
		data[i+1], data[end] = data[end], data[i+1]
		sendUpdate(conn, data, i+1, getStepDelay)
		pi := i + 1

		quicksortWithProgress(conn, data, start, pi-1, getStepDelay)
		quicksortWithProgress(conn, data, pi+1, end, getStepDelay)
	}
}

func mergeSortWithProgress(conn *websocket.Conn, data []int, start, end int, getStepDelay func() time.Duration) {
	if end-start < 2 {
		return
	}

	mid := start + (end-start)/2
	mergeSortWithProgress(conn, data, start, mid, getStepDelay)
	mergeSortWithProgress(conn, data, mid, end, getStepDelay)

	temp := make([]int, len(data))
	i, j, k := start, mid, start

	for i < mid && j < end {
		if data[i] < data[j] {
			temp[k] = data[i]
			i++
		} else {
			temp[k] = data[j]
			j++
		}
		k++
	}

	for i < mid {
		temp[k] = data[i]
		i++
		k++
	}
	for j < end {
		temp[k] = data[j]
		j++
		k++
	}

	for i = start; i < end; i++ {
		data[i] = temp[i]
		sendUpdate(conn, data, i, getStepDelay)
	}
	sendUpdate(conn, data, end-1, getStepDelay)
}

func heapSortWithProgress(conn *websocket.Conn, data []int, getStepDelay func() time.Duration) {
	n := len(data)

	for i := n/2 - 1; i >= 0; i-- {
		heapifyWithProgress(conn, data, n, i, getStepDelay)
	}

	for i := n - 1; i > 0; i-- {
		data[0], data[i] = data[i], data[0]
		sendUpdate(conn, data, i, getStepDelay)
		heapifyWithProgress(conn, data, i, 0, getStepDelay)
	}
}

func heapifyWithProgress(conn *websocket.Conn, data []int, n, i int, getStepDelay func() time.Duration) {
	largest := i
	l := 2*i + 1
	r := 2*i + 2

	if l < n && data[l] > data[largest] {
		largest = l
	}

	if r < n && data[r] > data[largest] {
		largest = r
	}
	if largest != i {
		data[i], data[largest] = data[largest], data[i]
		sendUpdate(conn, data, largest, getStepDelay)
		heapifyWithProgress(conn, data, n, largest, getStepDelay)
	}
}

func radixSortWithProgress(conn *websocket.Conn, data []int, getStepDelay func() time.Duration) {
	max := data[0]
	for _, value := range data[1:] {
		if value > max {
			max = value
		}
	}
	for exp := 1; max/exp > 0; exp *= 10 {
		n := len(data)
		output := make([]int, n)
		count := make([]int, 10)

		for i := 0; i < n; i++ {
			index := (data[i] / exp) % 10
			count[index]++
		}

		for i := 1; i < 10; i++ {
			count[i] += count[i-1]
		}

		for i := n - 1; i >= 0; i-- {
			index := (data[i] / exp) % 10
			output[count[index]-1] = data[i]
			count[index]--
		}

		for i := 0; i < n; i++ {
			data[i] = output[i]
			sendUpdate(conn, data, i, getStepDelay)
		}
	}
}

func bogoSortWithProgress(conn *websocket.Conn, data []int, getStepDelay func() time.Duration) {
	for !isSorted(conn, data, getStepDelay) {
		shuffle(data)
	}
}

func isSorted(conn *websocket.Conn, data []int, getStepDelay func() time.Duration) bool {
	for i := 1; i < len(data); i++ {
		sendUpdate(conn, data, i, getStepDelay)
		if data[i] < data[i-1] {
			return false
		}
	}
	return true
}

func shuffle(data []int) {
	rand.Shuffle(len(data), func(i, j int) {
		data[i], data[j] = data[j], data[i]
	})
}

func sendUpdate(conn *websocket.Conn, data []int, lastConsidered int, getStepDelay func() time.Duration) {
	update := SortMessage{
		Type:           "progress",
		Data:           data,
		LastConsidered: lastConsidered,
	}
	msg, _ := json.Marshal(update)
	conn.WriteMessage(websocket.TextMessage, msg)
	startTime := time.Now()
	targetDuration := getStepDelay()
	for time.Since(startTime) < targetDuration {
		// do nothing
	}
}
