package main

import (
	"github.com/TR-Projekt/zwischentoncloud/server"
	"github.com/TR-Projekt/zwischentoncloud/server/config"
	"github.com/TR-Projekt/zwischentoncloud/server/logger"
	"github.com/rs/zerolog/log"
)

func main() {

	logger.InitializeGlobalLogger("/var/log/zwischentoncloud/info.log", true)
	log.Info().Msg("Server startup.")

	conf := config.DefaultConfig()
	log.Info().Msg("Server configuration was initialized.")

	config.CheckForArguments()

	server := server.NewServer(conf)
	go server.Run(conf)
	log.Info().Msg("Server did start.")

	// wait forever
	// https://stackoverflow.com/questions/36419054/go-projects-main-goroutine-sleep-forever
	select {}
}
