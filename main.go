package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/caarlos0/env/v7"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	sunriseSunsetAPIURL = "https://api.meteo-concept.com/api/ephemeride/0"
	warningMessage      = "<!here> Ça va être tout noir"
)

type Ephemeride struct {
	City struct {
		Insee     string  `json:"insee"`
		Cp        string  `json:"cp"`
		Name      string  `json:"name"`
		Altitude  int     `json:"altitude"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"city"`
	E struct {
		Insee           string  `json:"insee"`
		Datetime        string  `json:"datetime"`
		Sunrise         string  `json:"sunrise"`
		Sunset          string  `json:"sunset"`
		DurationDay     string  `json:"duration_day"`
		MoonPhase       string  `json:"moon_phase"`
		Day             int     `json:"day"`
		DiffDurationDay int     `json:"diff_duration_day"`
		Latitude        float64 `json:"latitude"`
		Longitude       float64 `json:"longitude"`
		MoonAge         float64 `json:"moon_age"`
	} `json:"ephemeride"`
}

func getEphemeride(conf MeteoConf) (Ephemeride, error) {
	var e Ephemeride
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s?insee=%s", sunriseSunsetAPIURL, conf.Location),
		nil,
	)

	if err != nil {
		return e, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", conf.Token))
	req.Header.Add("Accept", "application/json")

	for tries := 0; tries < conf.Tries; tries++ {
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return e, err
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return e, err
		}

		if res.StatusCode > 500 {
			time.Sleep(time.Duration(10*(tries+1)) * time.Second)
			fmt.Println("waiting for API")
			continue
		}

		if res.StatusCode > 400 {
			return e, fmt.Errorf("credentials problems: %s", string(body))
		}

		if err := json.Unmarshal(body, &e); err != nil {
			return e, err
		}

		return e, nil
	}

	return e, fmt.Errorf("too many tries")
}

func getSunsetHours(conf MeteoConf) (int, int, error) {
	ephemeride, err := getEphemeride(conf)
	if err != nil {
		return 0, 0, err
	}

	sunset := strings.Split(ephemeride.E.Sunset, ":")

	sunsetHour, err := strconv.Atoi(sunset[0])
	if err != nil {
		return 0, 0, err
	}
	sunsetMinute, err := strconv.Atoi(sunset[1])
	if err != nil {
		return 0, 0, err
	}

	return sunsetHour, sunsetMinute, nil
}

func warnItsGonnaGetDark(conf SlackConf) error {
	warnMsg := []byte(fmt.Sprintf(`{"text": "%s", "channel": "%s"}`, warningMessage, conf.Channel))
	slackReq, err := http.NewRequest(
		http.MethodPost,
		conf.APIURL,
		bytes.NewReader(warnMsg),
	)

	if err != nil {
		return err
	}
	slackReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", conf.Token))
	slackReq.Header.Add("Content-Type", "application/json")

	fmt.Println("Ça va être tout noir")

	res, err := http.DefaultClient.Do(slackReq)
	if err != nil {
		return err
	}

	if res.StatusCode > 300 {
		fmt.Println(io.ReadAll(res.Body))
	}

	return nil
}

func launchWarning(conf Config) error {
	for {
		now := time.Now()
		l, err := time.LoadLocation("Europe/Paris")

		sunsetHour, sunsetMinute, err := getSunsetHours(conf.Meteo)
		if err != nil {
			return err
		}

		startupTime := time.Date(now.Year(), now.Month(), now.Day(), sunsetHour, sunsetMinute-1, 0, 0, l)
		fmt.Printf("Speaking at %s\n", startupTime.String())

		warnChan := make(chan error, 1)
		time.AfterFunc(startupTime.Sub(now), func() {
			warnChan <- warnItsGonnaGetDark(conf.Slack)
		})

		select {
		case err := <-warnChan:
			if err != nil {
				return err
			}
			tomorrow := time.Date(now.Year(), now.Month(), now.Day()+1, now.Hour(), 0, 0, 0, l)
			fmt.Printf("next time will be %s\n", tomorrow.String())
			time.Sleep(tomorrow.Sub(now))
			continue
		}
	}
}

func run() error {
	conf := Config{}
	if err := env.Parse(&conf); err != nil {
		return err
	}

	return launchWarning(conf)
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
