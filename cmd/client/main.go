package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"ecg_tool/pkg/detection"
	"ecg_tool/pkg/model"
	"ecg_tool/pkg/notification"

	"github.com/gorilla/websocket"
)

var (
	addr       = flag.String("addr", "localhost:8080", "server address")
	logDir     = flag.String("logdir", "logs", "directory to store log files")
	enableBeep = flag.Bool("beep", true, "enable beep sounds for alerts")
	logNormal  = flag.Bool("lognormal", true, "log normal readings (not just irregularities), defaults to true now")
	debug      = flag.Bool("debug", false, "enable debug mode for extra logging")
)

func main() {
	flag.Parse()
	log.SetFlags(0)

	detector := detection.NewDetector()

	logger, err := notification.NewLogger(*logDir)
	if err != nil {
		log.Fatalf("Failed to set up logging: %v", err)
	}
	defer logger.Close()

	// Set up beeper
	beeper := notification.NewBeeper(*enableBeep)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Connect to WebSocket server
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ecg"}
	log.Printf("Connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	log.Printf("Connected to ECG Sender")
	log.Printf("Logs will be stored in: %s", filepath.Join(*logDir))
	log.Printf("Beep alerts enabled: %v", *enableBeep)
	log.Printf("Logging all ECG data: enabled")

	if !*debug {
		fmt.Println("\n=== ECG Monitoring Started ===")
		fmt.Println("Legend: . = normal reading, ! = irregularity detected")
		fmt.Print("Monitoring: ")
	}

	done := make(chan struct{})
	messageCount := 0

	// Read messages from WebSocket
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}

			messageCount++

			var ecgData model.ECGData
			if err := json.Unmarshal(message, &ecgData); err != nil {
				log.Println("Error parsing ECG data:", err)
				continue
			}

			condition := detector.AnalyzeECGData(ecgData)

			if condition != nil {
				logger.LogCondition(condition)

				if !*debug {
					fmt.Print("!")
				}
				if *enableBeep {
					if *debug {
						fmt.Printf("\nAttempting to beep for %s (severity: %s)...\n",
							condition.Type, condition.Severity)
					}
					beeper.BeepForCondition(condition)
				}
			} else {
				logger.LogNormalReading(ecgData)

				if !*debug && messageCount%5 == 0 {
					fmt.Print(".")
				}
			}
			if !*debug && messageCount%50 == 0 {
				fmt.Print("\nMonitoring: ")
			}
		}
	}()

	<-interrupt
	fmt.Println("\n\nShutting down...")

	err = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Println("write close:", err)
	}

	select {
	case <-done:
	case <-time.After(time.Second):
	}
}
