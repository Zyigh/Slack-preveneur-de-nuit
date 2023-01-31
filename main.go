package main

import (
	"github.com/caarlos0/env/v7"
	"log"
	"preveneurdenuit/lib"
	"preveneurdenuit/lib/config"
	"preveneurdenuit/lib/notifiers"
	"preveneurdenuit/lib/providers"
)

func run() error {
	conf := config.Config{}
	if err := env.Parse(&conf); err != nil {
		return err
	}

	app := lib.NewApp(providers.NewMeteoConcept(conf.Meteo))
	app = app.WithNotifier(notifiers.NewSlack(conf.Slack))

	return app.Start()
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
