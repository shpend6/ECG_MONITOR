package notification

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"

	"ecg_tool/pkg/model"
)

type Beeper struct {
	enabled bool
}

func NewBeeper(enabled bool) *Beeper {
	return &Beeper{
		enabled: enabled,
	}
}

func (b *Beeper) BeepForCondition(condition *model.HeartCondition) {
	if !b.enabled || condition == nil {
		return
	}

	err := b.PlaySystemSound()
	if err != nil {
		log.Printf("Warning: System sound failed: %v, falling back to console beep", err)

		switch condition.Type {
		case "tachycardia":
			b.playFastBeep(condition.Severity)
		case "bradycardia":
			b.playSlowBeep(condition.Severity)
		case "arrhythmia":
			b.playIrregularBeep(condition.Severity)
		}
	}
}

func (b *Beeper) playFastBeep(severity string) {
	count := 1
	if severity == "moderate" {
		count = 2
	} else if severity == "severe" {
		count = 3
	}

	for i := 0; i < count; i++ {
		fmt.Print("\a")
		time.Sleep(100 * time.Millisecond)
	}
}

func (b *Beeper) playSlowBeep(severity string) {
	count := 1
	if severity == "moderate" {
		count = 2
	} else if severity == "severe" {
		count = 3
	}

	for i := 0; i < count; i++ {
		fmt.Print("\a")
		time.Sleep(300 * time.Millisecond)
	}
}

func (b *Beeper) playIrregularBeep(severity string) {
	if severity == "severe" {
		fmt.Print("\a")
		time.Sleep(100 * time.Millisecond)
		fmt.Print("\a")
		time.Sleep(300 * time.Millisecond)
		fmt.Print("\a")
		time.Sleep(100 * time.Millisecond)
		fmt.Print("\a")
		return
	}

	fmt.Print("\a")
	time.Sleep(200 * time.Millisecond)
	fmt.Print("\a")
}

func (b *Beeper) PlaySystemSound() error {
	if !b.enabled {
		return nil
	}

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		// Try several sound methods for Linux
		if _, err := exec.LookPath("paplay"); err == nil {
			// Try Pulse Audio first (most common on modern systems)
			soundFile := "/usr/share/sounds/freedesktop/stereo/bell.oga"
			if _, err := os.Stat(soundFile); os.IsNotExist(err) {
				soundFile = "/usr/share/sounds/freedesktop/stereo/complete.oga"
				if _, err := os.Stat(soundFile); os.IsNotExist(err) {
					// Try alternative sound
					soundFile = "/usr/share/sounds/sound-icons/bell.wav"
				}
			}
			cmd = exec.Command("paplay", soundFile)
		}
	case "darwin":
		// macOS
		cmd = exec.Command("afplay", "/System/Library/Sounds/Ping.aiff")
	case "windows":
		// Windows - using PowerShell to play a sound
		cmd = exec.Command("powershell", "-c", "[console]::beep(1000,500)")
	}

	if cmd == nil {
		return fmt.Errorf("no suitable audio player found for OS: %s", runtime.GOOS)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
