package rrmultivoice

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var fileDirIndex = make(map[string]int) // Contains last used folder index within user supplied sound packs root directory
var lastPlayedFile string               // Contains currently played (or last played) sound's filename

func scanInUseFiles() {
	fds, err := os.ReadDir(fdPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			updatePID()
		} else {
			log.Println("Failed to read", fdPath, "dir:", err)
			time.Sleep(retryDuration)
		}
		scanInUseFiles()
		return
	}

	fileName := ""
	fileDetected := false
	for _, fd := range fds {
		link, err := os.Readlink(fdPath + "/" + fd.Name())
		if err != nil {
			log.Println("Failed to update RoboController PID:", err)
			continue
		}

		if strings.HasPrefix(link, soundsDir) && !strings.HasSuffix(link, " (deleted)") {
			fileName = filepath.Base(link)
			fileDetected = true
			break
		}
	}

	if fileDetected {
		if fileName != lastPlayedFile {
			lastPlayedFile = fileName
			log.Println("Replacing file:", fileName)
			linklastPlayedFile()
		}
	} else if lastPlayedFile != "" {
		lastPlayedFile = ""
	}
}

func linklastPlayedFile() {
	lastIndex := fileDirIndex[lastPlayedFile] // If not found = default = 0

	// Scan user supplied sound dirs root directory
	soundPacks, err := os.ReadDir(packsDir)
	if err != nil {
		log.Println("Failed to list dirs in user-provided packs dir:", err)
		return
	}

	// Collect list of dirs inside user supplied sound dirs root directory
	detectedDirs := make([]string, 0, len(soundPacks))
	for _, soundPack := range soundPacks {
		if !soundPack.IsDir() {
			continue
		}
		detectedDirs = append(detectedDirs, soundPack.Name())
	}
	if len(detectedDirs) == 0 {
		log.Println("No user provided sound packs detected in", packsDir)
		return
	}

	// Increase index by 1
	newDirIndex := lastIndex + 1

	// Search for a file in a folder at a new index
	for {

		// Check if we performed a full cycle, yet did not find a required file in new directories
		if newDirIndex == lastIndex {
			break
		}

		// Check if index is higher than detect dirs count-1
		if len(detectedDirs)-1 < newDirIndex {
			newDirIndex = 0
		}

		// Check if file exists at folder with new index
		if _, err := os.Stat(packsDir + detectedDirs[newDirIndex]); errors.Is(err, os.ErrNotExist) {
			newDirIndex++
			continue
		}

		hardLinkSrc := packsDir + detectedDirs[newDirIndex] + "/" + lastPlayedFile
		hardLinkDst := soundsDir + lastPlayedFile

		// Delete old file/hardlink
		err = os.Remove(hardLinkDst)
		if err != nil {
			log.Println("Failed to delete file/hardlink:", err)
			// Do not "return" here. File might not even exist, right?
		}

		// Create hard link
		err = os.Link(hardLinkSrc, hardLinkDst)
		if err != nil {
			log.Println("Failed to create hardlink from", hardLinkSrc, "to", hardLinkDst+":", err)
			return
		}
		break
	}

	fileDirIndex[lastPlayedFile] = newDirIndex
}
