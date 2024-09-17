package main

import (
	"github.com/root4loot/goutils/log"
)

func main() {

	log.Init("Myapp")
	if log.IsOutputPiped() {
		log.Notify(log.PipedOutputNotification)
	}

	log.SetLevel(log.TraceLevel)

	log.Info("Some info")
	log.Result("Some result")
	log.Warn("Some warning")

	newLogger := log.NewLogger("myApp2")
	newLogger.Info("Some info")

}
