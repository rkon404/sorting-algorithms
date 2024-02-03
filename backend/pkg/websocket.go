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

func NewHandler() *Handler {
	return &Handler{StepDelay: 50 * time.Millisecond}
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
					bubbleSortWithProgress(conn, sortMsg.Data, handler.GetStepDelay)
				}
				if sortMsg.Algorithm == "quick" {
					go quicksortWithProgress(conn, sortMsg.Data, 0, len(sortMsg.Data)-1, handler.GetStepDelay)
				}
				if sortMsg.Algorithm == "merge" {
					mergeSortWithProgress(conn, sortMsg.Data, 0, len(sortMsg.Data), handler.GetStepDelay)
				}
				if sortMsg.Algorithm == "heap" {
					heapSortWithProgress(conn, sortMsg.Data, handler.GetStepDelay)
				}
				if sortMsg.Algorithm == "radix" {
					radixSortWithProgress(conn, sortMsg.Data, handler.GetStepDelay)
				}
				if sortMsg.Algorithm == "bongo" {
					bongoSortWithProgress(conn, sortMsg.Data, handler.GetStepDelay)
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
