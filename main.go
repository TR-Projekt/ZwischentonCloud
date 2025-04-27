package main

import (
	servertools "github.com/Festivals-App/festivals-server-tools"
	"github.com/TR-Projekt/zwischentoncloud/server"
	"github.com/TR-Projekt/zwischentoncloud/server/config"
	"github.com/rs/zerolog/log"
)

func main() {

	log.Info().Msg("Server startup.")

	root := servertools.ContainerPathArgument()
	configFilePath := root + "/etc/zwischentoncloud.conf"

	conf := config.ParseConfig(configFilePath)
	log.Info().Msg("Server configuration was initialized")

	servertools.InitializeGlobalLogger(conf.InfoLog, true)
	log.Info().Msg("Logger initialized")

	server := server.NewServer(conf)
	go server.Run(conf)
	log.Info().Msg("Server did start")

	// wait forever
	// https://stackoverflow.com/questions/36419054/go-projects-main-goroutine-sleep-forever
	select {}
}
