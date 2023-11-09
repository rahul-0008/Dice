package main

import (
	"flag"
	"log"

	"github.com/DiceDB/Dice/config"
	"github.com/DiceDB/Dice/server"
)

func setupFlags() {
	flag.StringVar(&config.Host, "Host", "0.0.0.0", "The host for the dice db")
	flag.IntVar(&config.Port, "Port", 7379, "Port for thre dice server")

	flag.Parse()
}

func main() {
	setupFlags()
	log.Println("Rolling the Dice 🎲")
	server.RunSyncTCPServer()

}
