package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"ecg_tool/pkg/model"
	"ecg_tool/pkg/notification"
)

func main() {
	log.Println("Starting sound test utility")

	// Sound generation mode
	mode := flag.String("mode", "all", "Sound test mode: tachycardia, bradycardia, arrhythmia, system, all")
	severity := flag.String("severity", "severe", "Severity level: mild, moderate, severe")
	flag.Parse()

	beeper := notification.NewBeeper(true)

	if *mode == "all" || *mode == "system" {
		log.Println("Testing system sound...")
		err := beeper.PlaySystemSound()
		if err != nil {
			log.Printf("System sound failed: %v", err)
		} else {
			log.Println("System sound played successfully")
		}
		time.Sleep(1 * time.Second)
	}

	if *mode == "all" || *mode == "tachycardia" {
		log.Printf("Testing tachycardia beep (%s)...", *severity)
		condition := &model.HeartCondition{
			Type:        "tachycardia",
			Severity:    *severity,
			HeartRate:   120,
			Timestamp:   time.Now(),
			Description: "Test tachycardia",
		}
		beeper.BeepForCondition(condition)
		time.Sleep(1 * time.Second)
	}

	if *mode == "all" || *mode == "bradycardia" {
		log.Printf("Testing bradycardia beep (%s)...", *severity)
		condition := &model.HeartCondition{
			Type:        "bradycardia",
			Severity:    *severity,
			HeartRate:   45,
			Timestamp:   time.Now(),
			Description: "Test bradycardia",
		}
		beeper.BeepForCondition(condition)
		time.Sleep(1 * time.Second)
	}

	if *mode == "all" || *mode == "arrhythmia" {
		log.Printf("Testing arrhythmia beep (%s)...", *severity)
		condition := &model.HeartCondition{
			Type:        "arrhythmia",
			Severity:    *severity,
			HeartRate:   70,
			Timestamp:   time.Now(),
			Description: "Test arrhythmia",
		}
		beeper.BeepForCondition(condition)
	}

	fmt.Println("\nTest completed. Did you hear any sounds?")

}
