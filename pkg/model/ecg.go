package model

import (
	"time"
)

type ECGData struct {
	Timestamp     time.Time `json:"timestamp"`
	HeartRate     int       `json:"heart_rate"`
	RRInterval    float64   `json:"rr_interval"`
	QTInterval    float64   `json:"qt_interval"`
	PRInterval    float64   `json:"pr_interval"`
	QRSInterval   float64   `json:"qrs_interval"`
	SignalQuality float64   `json:"signal_quality"`
}

type HeartCondition struct {
	Type        string    `json:"type"`
	Severity    string    `json:"severity"`
	HeartRate   int       `json:"heart_rate"`
	Timestamp   time.Time `json:"timestamp"`
	Description string    `json:"description"`
}

const (
	TachycardiaThreshold = 100
	BradycardiaThreshold = 60

	RRVariationArrhythmiaThreshold = 0.15
)
