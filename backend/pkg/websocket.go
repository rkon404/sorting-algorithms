package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Handler struct {
	StepDelay time.Duration
	mu        sync.Mutex
}

func (h *Handler) SetStepDelay(delay float32) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.StepDelay = time.Duration(delay) * time.Millisecond
}

func (h *Handler) GetStepDelay() time.Duration {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.StepDelay
}

func NewHandler() *Handler {
	return &Handler{StepDelay: 100 * time.Millisecond}
}

type SortMessage struct {
	Type           string `json:"type"`
	Algorithm      string `json:"algorithm"`
	Data           []int  `json:"data"`
	LastConsidered int    `json:"lastConsidered,omitempty"`
}

type ConfirmationMessage struct {
	Type    string  `json:"type"`
	Success bool    `json:"success"`
	Delay   float32 `json:"delay"`
}

type StepDelayMessage struct {
	Type string  `json:"type"`
	Data float32 `json:"data"`
}

func (h *Handler) bubbleSortWithProgress(conn *websocket.Conn, data []int) {
	n := len(data)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if data[j] > data[j+1] {

				data[j], data[j+1] = data[j+1], data[j]

				update := SortMessage{
					Type:           "progress",
					Data:           data,
					LastConsidered: j + 1,
				}

				msg, _ := json.Marshal(update)
				conn.WriteMessage(websocket.TextMessage, msg)

				startTime := time.Now()
				targetDuration := h.GetStepDelay()
				for time.Since(startTime) < targetDuration {
					// do nothing
				}
			}
		}
	}
}

func (h *Handler) quicksortWithProgress(conn *websocket.Conn, data []int, start, end int) {
	if start < end {
		pi := h.partition(conn, data, start, end)

		h.quicksortWithProgress(conn, data, start, pi-1)
		h.quicksortWithProgress(conn, data, pi+1, end)
	}
}

func (h *Handler) partition(conn *websocket.Conn, data []int, start, end int) int {
	pivot := data[end]
	i := start - 1

	for j := start; j < end; j++ {
		if data[j] <= pivot {
			i++
			data[i], data[j] = data[j], data[i]

			h.sendUpdate(conn, data, j)
		}
	}
	data[i+1], data[end] = data[end], data[i+1]

	h.sendUpdate(conn, data, i+1)

	return i + 1
}

func (h *Handler) sendUpdate(conn *websocket.Conn, data []int, lastConsidered int) {
	update := SortMessage{
		Type:           "progress",
		Data:           data,
		LastConsidered: lastConsidered,
	}

	msg, _ := json.Marshal(update)
	conn.WriteMessage(websocket.TextMessage, msg)
	startTime := time.Now()
	targetDuration := h.GetStepDelay()
	for time.Since(startTime) < targetDuration {
		// do nothing
	}
}

func (h *Handler) mergeSortWithProgress(conn *websocket.Conn, data []int, start, end int) {
	if end-start < 2 {
		return
	}

	tempData := make([]int, len(data))
	mid := start + (end-start)/2
	h.mergeSortWithProgress(conn, data, start, mid)
	h.mergeSortWithProgress(conn, data, mid, end)
	h.merge(conn, data, tempData, start, mid, end)
}

func (h *Handler) merge(conn *websocket.Conn, data, temp []int, start, mid, end int) {
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
	}

	h.sendMergeUpdate(conn, data, k-1)
}

func (h *Handler) sendMergeUpdate(conn *websocket.Conn, data []int, lastConsidered int) {
	update := SortMessage{
		Type:           "progress",
		Data:           data,
		LastConsidered: lastConsidered,
	}
	msg, _ := json.Marshal(update)
	conn.WriteMessage(websocket.TextMessage, msg)
	startTime := time.Now()
	targetDuration := h.GetStepDelay()
	for time.Since(startTime) < targetDuration {
		// do nothing
	}
}

func ServeWs(handler *Handler, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		var baseMsg map[string]interface{}
		if err := json.Unmarshal(message, &baseMsg); err != nil {
			log.Println("Message unmarshal error:", err)
			continue
		}

		switch baseMsg["type"] {
		case "sort":
			var sortMsg SortMessage
			if err := json.Unmarshal(message, &sortMsg); err != nil {
				log.Println("Sort message unmarshal error:", err)
				continue
			}
			go func() {
				if sortMsg.Algorithm == "bubble" {
					handler.bubbleSortWithProgress(conn, sortMsg.Data)
				}
				if sortMsg.Algorithm == "quick" {
					go handler.quicksortWithProgress(conn, sortMsg.Data, 0, len(sortMsg.Data)-1)
				}
				if sortMsg.Algorithm == "merge" {
					handler.mergeSortWithProgress(conn, sortMsg.Data, 0, len(sortMsg.Data))

				}
			}()
		case "stepDelay":
			var delayMsg StepDelayMessage
			if err := json.Unmarshal(message, &delayMsg); err != nil {
				log.Println("Step delay message unmarshal error:", err)
				continue
			}
			handler.SetStepDelay(delayMsg.Data)

			confirmation := ConfirmationMessage{
				Type:    "stepDelayConfirmation",
				Success: true,
				Delay:   delayMsg.Data,
			}
			confMsg, err := json.Marshal(confirmation)
			if err != nil {
				log.Println("Error marshaling confirmation message:", err)
				continue
			}
			if err := conn.WriteMessage(websocket.TextMessage, confMsg); err != nil {
				log.Println("Error sending confirmation message:", err)
			}
		}
	}
}
