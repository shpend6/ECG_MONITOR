package main

import (
	"encoding/json"
	"flag"
	"log"
	"math/rand"
	"net/http"
	"time"

	"ecg_tool/pkg/model"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var (
	addr              = flag.String("addr", ":8080", "http service address")
	sendInterval      = flag.Duration("interval", 1*time.Second, "interval between sending data (e.g. 1s, 500ms)")
	simulateIrregular = flag.Bool("irregular", true, "simulate irregular heartbeats")
)

func main() {
	flag.Parse()
	log.SetFlags(0)

	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/ecg", handleECGConnection)

	log.Printf("ECG Sender starting on %s\n", *addr)
	log.Printf("Send interval: %v\n", *sendInterval)
	log.Printf("Simulate irregular heartbeats: %v\n", *simulateIrregular)

	log.Fatal(http.ListenAndServe(*addr, nil))
}

func handleECGConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	log.Println("Client connected")

	ticker := time.NewTicker(*sendInterval)
	defer ticker.Stop()

	heartbeats := 0

	for range ticker.C {
		ecgData := generateECGData(heartbeats)
		heartbeats++

		jsonData, err := json.Marshal(ecgData)
		if err != nil {
			log.Println("Error marshaling data:", err)
			continue
		}

		if err := conn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
			log.Println("Error sending data:", err)
			break
		}

		log.Printf("Sent data: Heart rate=%d, RR=%.2fs", ecgData.HeartRate, ecgData.RRInterval)
	}
}

func generateECGData(heartbeats int) model.ECGData {
	var heartRate int
	var rrInterval float64

	baseHeartRate := 70 + rand.Intn(15)

	baseRRInterval := 60.0 / float64(baseHeartRate)

	if *simulateIrregular && (heartbeats%30 == 0) {
		switch rand.Intn(3) {
		case 0: //tachycardia
			heartRate = 100 + rand.Intn(40)
			rrInterval = 60.0 / float64(heartRate)
		case 1: //bradycardia
			heartRate = 40 + rand.Intn(20)
			rrInterval = 60.0 / float64(heartRate)
		case 2: //normal
			heartRate = baseHeartRate
			rrInterval = baseRRInterval * (1.0 + (0.2 * (rand.Float64() - 0.5)))
		}
	} else {
		heartRate = baseHeartRate + rand.Intn(5) - 2
		rrInterval = baseRRInterval * (1.0 + (0.05 * (rand.Float64() - 0.5)))
	}

	qtInterval := 0.35 + (0.02 * rand.Float64())
	prInterval := 0.15 + (0.05 * rand.Float64())
	qrsInterval := 0.08 + (0.02 * rand.Float64())

	signalQuality := 0.9 + (0.1 * rand.Float64())
	if rand.Intn(100) < 5 {
		signalQuality = 0.5 + (0.4 * rand.Float64())
	}

	return model.ECGData{
		Timestamp:     time.Now(),
		HeartRate:     heartRate,
		RRInterval:    rrInterval,
		QTInterval:    qtInterval,
		PRInterval:    prInterval,
		QRSInterval:   qrsInterval,
		SignalQuality: signalQuality,
	}
}
