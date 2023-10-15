package main

import (
	"flag"
	"log"
	"rrmultivoice/pkg/rrmultivoice"
	"time"
)

var (
	flagInterval  = flag.Duration("interval", 300*time.Millisecond, "Interval for scanning fds.")
	flagSoundsDir = flag.String("soundsdir", "/opt/rockrobo/resources/sounds/en/", "Original directory from which robot plays voice-lines.")
	flagPacksDir  = flag.String("packsdir", "/opt/rockrobo/resources/sounds/en/additional_sounds", "Directory of additional sound directories")
)

func main() {
	flag.Parse()

	// Ensure interval is not abnormally low
	if *flagInterval < 10*time.Millisecond {
		log.Fatalln("Interval must be at least 10ms. Ideally within 200m-1000ms range")
	}

	rrmultivoice.Start(*flagInterval, *flagSoundsDir, *flagPacksDir)
}
