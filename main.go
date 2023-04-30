package main

import (
	"github.com/lent0s/provider-status/cmd"
	"log"
	"os"
)

func main() {

	logFile, err := os.OpenFile("status.log",
		os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal("access to log file denied")
	}
	log.SetOutput(logFile)

	cmd.ReadConfig()
	cmd.ServerUp()
}
