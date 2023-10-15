package rrmultivoice

import (
	"log"
	"os"
	"strconv"
	"time"
)

var (
	controllerPID string
	fdPath        string
)

func updatePID() {
	files, err := os.ReadDir("/proc")
	if err != nil {
		log.Println("Failed to read PID (/proc dir):", err)
		time.Sleep(retryDuration)
		updatePID()
		return
	}

	for _, file := range files {
		if !file.IsDir() {
			continue // Not a dir
		}

		pid := file.Name()
		_, err := strconv.Atoi(pid)
		if err != nil {
			continue // Not a process dir
		}

		cmdline, err := os.ReadFile("/proc/" + pid + "/cmdline")
		if err != nil {
			log.Println("Failed to read PID (/proc/"+pid+"/cmdline dir):", err)
			time.Sleep(retryDuration)
			updatePID()
			return
		}

		if string(cmdline) == "RoboController\x00" {
			log.Println("Found RoboController PID:", pid)
			controllerPID = pid
			fdPath = "/proc/" + controllerPID + "/fd"
			return
		}
	}

	log.Println("Failed to read PID: RoboController PID not found")
	time.Sleep(retryDuration)
	updatePID()
}
