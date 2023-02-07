package lib

import (
	"fmt"
	"log"
	"time"

	"github.com/Zyigh/Slack-preveneur-de-nuit/lib/notifiers"
	"github.com/Zyigh/Slack-preveneur-de-nuit/lib/utils"
)

const warningMessage = "<!here> Ça va être tout noir !"

type Notifier interface {
	Notify(msg string) error
}

type SunsetHourAndMinuter interface {
	SunsetHourAndMinute() (int, int, error)
}

type App struct {
	provider  SunsetHourAndMinuter
	notifiers []Notifier
}

func NewApp(provider SunsetHourAndMinuter) App {
	return App{
		provider: provider,
	}
}

func (a App) WithNotifier(notifiers ...Notifier) App {
	a.notifiers = append(a.notifiers, notifiers...)

	return a
}

func (a App) Start() error {
	if len(a.notifiers) == 0 {
		a.notifiers = append(a.notifiers, notifiers.Default{})
	}

	for {
		hour, minute, err := a.provider.SunsetHourAndMinute()
		if err != nil {
			return fmt.Errorf("run: can't get sunset: %w", err)
		}
		now := time.Now()
		l, err := time.LoadLocation("Europe/Paris")
		startupTime := time.Date(now.Year(), now.Month(), now.Day(), hour, minute-1, 0, 0, l)

		if startupTime.After(now) {
			log.Printf("Speaking at %s\n", startupTime.String())

			warnChan := make(chan error, 1)
			time.AfterFunc(startupTime.Sub(now), func() {
				errs := make([]error, 0, len(a.notifiers))
				for _, notifier := range a.notifiers {
					if err := notifier.Notify(warningMessage); err != nil {
						errs = append(errs, err)
					}
				}

				if len(errs) > 0 {
					warnChan <- utils.Reduce(errs, func(acc error, err error) error {
						return fmt.Errorf("%s\n%w", acc, err)
					}, fmt.Errorf(""))
					return
				}

				warnChan <- nil
			})

			select {
			case err := <-warnChan:
				if err != nil {
					return err
				}
			}
		}

		tomorrow := time.Date(now.Year(), now.Month(), now.Day()+1, 8, 0, 0, 0, l)
		log.Printf("next time will be %s\n", tomorrow.String())
		time.Sleep(tomorrow.Sub(now))
	}
}
