package main

import (
	"github.com/ildomm/eskalationszeit/zeitarbeiter/config"
	"github.com/ildomm/eskalationszeit/zeitarbeiter/databases"
	"github.com/ildomm/eskalationszeit/zeitarbeiter/logic"
	"github.com/ildomm/eskalationszeit/zeitarbeiter/utils"
	"log"
	"time"
)

func main() {
	config.Setup()
	utils.SetupLogger("zeitarbeiter.log")

	log.Println("Starting worker Zeitarbeiter")
	utils.SignalNotify()

	databases.Setup()

	for true {
		logic.UpdatePrices()
		log.Printf("Next refresh in %d seconds", config.App.Runtime.RefreshSeconds)

		time.Sleep(time.Duration(config.App.Runtime.RefreshSeconds) * time.Second )
	}
}