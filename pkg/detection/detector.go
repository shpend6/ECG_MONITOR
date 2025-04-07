package detection

import (
	"math"
	"time"

	"ecg_tool/pkg/model"
)

type Detector struct {
	lastRRInterval float64
	lastTimestamp  time.Time
}

func NewDetector() *Detector {
	return &Detector{}
}

func (d *Detector) AnalyzeECGData(data model.ECGData) *model.HeartCondition {
	var condition *model.HeartCondition

	// Check for tachycardia (high heart rate)
	if data.HeartRate > model.TachycardiaThreshold {
		severity := "mild"
		if data.HeartRate > model.TachycardiaThreshold+20 {
			severity = "moderate"
		}
		if data.HeartRate > model.TachycardiaThreshold+40 {
			severity = "severe"
		}

		condition = &model.HeartCondition{
			Type:        "tachycardia",
			Severity:    severity,
			HeartRate:   data.HeartRate,
			Timestamp:   data.Timestamp,
			Description: "Abnormally high heart rate detected",
		}
	}

	if data.HeartRate < model.BradycardiaThreshold {
		severity := "mild"
		if data.HeartRate < model.BradycardiaThreshold-10 {
			severity = "moderate"
		}
		if data.HeartRate < model.BradycardiaThreshold-20 {
			severity = "severe"
		}

		condition = &model.HeartCondition{
			Type:        "bradycardia",
			Severity:    severity,
			HeartRate:   data.HeartRate,
			Timestamp:   data.Timestamp,
			Description: "Abnormally low heart rate detected",
		}
	}

	if d.lastRRInterval > 0 && !d.lastTimestamp.IsZero() {
		variation := math.Abs(data.RRInterval-d.lastRRInterval) / d.lastRRInterval

		if variation > model.RRVariationArrhythmiaThreshold {
			severity := "mild"
			if variation > model.RRVariationArrhythmiaThreshold*1.5 {
				severity = "moderate"
			}
			if variation > model.RRVariationArrhythmiaThreshold*2 {
				severity = "severe"
			}

			condition = &model.HeartCondition{
				Type:        "arrhythmia",
				Severity:    severity,
				HeartRate:   data.HeartRate,
				Timestamp:   data.Timestamp,
				Description: "Irregular heart rhythm detected",
			}
		}
	}

	// Update last values for next comparison
	d.lastRRInterval = data.RRInterval
	d.lastTimestamp = data.Timestamp

	return condition
}
