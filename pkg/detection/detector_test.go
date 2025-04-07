package detection

import (
	"testing"
	"time"

	"ecg_tool/pkg/model"
)

func TestDetectTachycardia(t *testing.T) {
	detector := NewDetector()

	data := model.ECGData{
		Timestamp:     time.Now(),
		HeartRate:     120,
		RRInterval:    60.0 / 120.0,
		QTInterval:    0.35,
		PRInterval:    0.15,
		QRSInterval:   0.08,
		SignalQuality: 0.95,
	}

	condition := detector.AnalyzeECGData(data)

	if condition == nil {
		t.Error("Expected tachycardia condition, got nil")
		return
	}

	if condition.Type != "tachycardia" {
		t.Errorf("Expected tachycardia, got %s", condition.Type)
	}
}

func TestDetectBradycardia(t *testing.T) {
	detector := NewDetector()

	data := model.ECGData{
		Timestamp:     time.Now(),
		HeartRate:     45,
		RRInterval:    60.0 / 45.0,
		QTInterval:    0.35,
		PRInterval:    0.15,
		QRSInterval:   0.08,
		SignalQuality: 0.95,
	}

	condition := detector.AnalyzeECGData(data)

	if condition == nil {
		t.Error("Expected bradycardia condition, got nil")
		return
	}

	if condition.Type != "bradycardia" {
		t.Errorf("Expected bradycardia, got %s", condition.Type)
	}
}

func TestDetectArrhythmia(t *testing.T) {
	detector := NewDetector()

	// First heartbeat (normal)
	data1 := model.ECGData{
		Timestamp:     time.Now(),
		HeartRate:     70,
		RRInterval:    0.85,
		QTInterval:    0.35,
		PRInterval:    0.15,
		QRSInterval:   0.08,
		SignalQuality: 0.95,
	}

	// Should be no condition for the first beat
	condition1 := detector.AnalyzeECGData(data1)
	if condition1 != nil {
		t.Errorf("Expected nil condition for first normal heartbeat, got %s", condition1.Type)
	}

	// Second heartbeat (significant RR interval change)
	data2 := model.ECGData{
		Timestamp:     time.Now().Add(time.Second),
		HeartRate:     70,
		RRInterval:    1.1, // >15% change from 0.85
		QTInterval:    0.35,
		PRInterval:    0.15,
		QRSInterval:   0.08,
		SignalQuality: 0.95,
	}

	// Should detect arrhythmia
	condition2 := detector.AnalyzeECGData(data2)
	if condition2 == nil {
		t.Error("Expected arrhythmia condition, got nil")
		return
	}

	if condition2.Type != "arrhythmia" {
		t.Errorf("Expected arrhythmia, got %s", condition2.Type)
	}
}

func TestNormalHeartbeat(t *testing.T) {
	detector := NewDetector()

	data := model.ECGData{
		Timestamp:     time.Now(),
		HeartRate:     70,
		RRInterval:    60.0 / 70.0,
		QTInterval:    0.35,
		PRInterval:    0.15,
		QRSInterval:   0.08,
		SignalQuality: 0.95,
	}

	condition := detector.AnalyzeECGData(data)

	if condition != nil {
		t.Errorf("Expected nil condition for normal heartbeat, got %s", condition.Type)
	}
}
