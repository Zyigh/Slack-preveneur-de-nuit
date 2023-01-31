package notifiers

import (
	"log"
)

type Default struct{}

type VerboseNotifier struct {
	Default
}

func (d Default) Notify(msg string) error {
	log.Println(msg)
	return nil
}
