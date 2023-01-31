package notifiers

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/Zyigh/Slack-preveneur-de-nuit/lib/config"
)

const slackMsgFormat = `{"text": "%s", "channel": "%s"}`

type Slack struct {
	apiURL  string
	token   string
	channel string
}

func NewSlack(conf config.SlackConf) Slack {
	return Slack{
		apiURL:  conf.APIURL,
		token:   conf.Token,
		channel: conf.Channel,
	}
}

func (s Slack) Notify(msg string) error {
	warnMsg := []byte(fmt.Sprintf(slackMsgFormat, msg, s.channel))
	slackReq, err := http.NewRequest(
		http.MethodPost,
		s.apiURL,
		bytes.NewReader(warnMsg),
	)

	if err != nil {
		return err
	}
	slackReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.token))
	slackReq.Header.Add("Content-Type", "application/json")

	log.Println("Sending message to Slack API")

	res, err := http.DefaultClient.Do(slackReq)
	if err != nil {
		return err
	}

	if res.StatusCode >= 300 {
		log.Println(io.ReadAll(res.Body))
	}

	return nil
}
