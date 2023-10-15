package rrmultivoice

import (
	"strings"
	"time"
)

const retryDuration = 5 * time.Second

var (
	scanInterval time.Duration
	soundsDir    string
	packsDir     string
)

func Start(i time.Duration, sd, pd string) {
	// Make them available globally
	scanInterval = i
	soundsDir = sd
	packsDir = pd

	// Ensure ends with slash
	if !strings.HasSuffix(soundsDir, "/") {
		soundsDir += "/"
	}
	if !strings.HasSuffix(packsDir, "/") {
		packsDir += "/"
	}

	updatePID()
	for {
		scanInUseFiles()
		time.Sleep(scanInterval)
	}
}
