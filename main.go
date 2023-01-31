package main

import (
	"log"

	"github.com/Zyigh/Slack-preveneur-de-nuit/lib"
	"github.com/Zyigh/Slack-preveneur-de-nuit/lib/config"
	"github.com/Zyigh/Slack-preveneur-de-nuit/lib/notifiers"
	"github.com/Zyigh/Slack-preveneur-de-nuit/lib/providers"
	"github.com/caarlos0/env/v7"
)

func run() error {
	conf := config.Config{}
	if err := env.Parse(&conf); err != nil {
		return err
	}

	app := lib.NewApp(providers.NewMeteoConcept(conf.Meteo)).
		WithNotifier(notifiers.NewSlack(conf.Slack))

	return app.Start()
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
