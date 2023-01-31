package providers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"preveneurdenuit/lib/config"
	"strconv"
	"strings"
	"time"
)

const meteoConceptAPIURL = "https://api.meteo-concept.com/api/ephemeride/0"

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

type MeteoConcept struct {
	apiURL   string
	token    string
	location string
	maxTries int
}

func NewMeteoConcept(conf config.MeteoConf) MeteoConcept {
	return MeteoConcept{
		token:    conf.Token,
		location: conf.Location,
		maxTries: conf.Tries,
		apiURL:   meteoConceptAPIURL,
	}
}

func (m MeteoConcept) getEphemeride() (Ephemeride, error) {
	var e Ephemeride

	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s?insee=%s", m.apiURL, m.location),
		nil,
	)

	if err != nil {
		return e, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", m.token))
	req.Header.Add("Accept", "application/json")

	for tries := 0; tries < m.maxTries; tries++ {
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

func (m MeteoConcept) SunsetHourAndMinute() (int, int, error) {
	ephemeride, err := m.getEphemeride()
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
