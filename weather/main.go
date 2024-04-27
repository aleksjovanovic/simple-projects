package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

type Weather struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`
	Current struct {
		TempC     float64 `json:"temp_c"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`
	Forecast struct {
		Forecastday []struct {
			Hour []struct {
				TimeEpoch int64   `json:"time_epoch"`
				TempC     float64 `json:"temp_c"`
				Condition struct {
					Text string `json:"text"`
				} `json:"condition"`
				ChanceOfRain float64 `json:"chance_of_rain"`
			} `json:"hour"`
		}
	} `json:"forecast"`
}

func main() {
	res, err := http.Get("http://api.weatherapi.com/v1/forecast.json?key=5aacad18601640e4875160518242704&q=Nis,Serbia&days=1&aqi=yes&alerts=no")
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Println("Weather API not available.")
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}
	var weather Weather
	err = json.Unmarshal(body, &weather)
	if err != nil {
		log.Println(err)
	}
	location, current, hours := weather.Location, weather.Current, weather.Forecast.Forecastday[0].Hour

	log.Printf(
		"%s,%s: %.0fC, %s\n",
		location.Name,
		location.Country,
		current.TempC,
		current.Condition.Text,
	)

	for _, hour := range hours {
		date := time.Unix(hour.TimeEpoch, 0)
		if date.Before(time.Now()) {
			continue
		}
		log.Printf(
			"%s - %.0fC, %.0f%%, %s\n",
			date.Format("15:04"),
			hour.TempC,
			hour.ChanceOfRain,
			hour.Condition.Text,
		)
	}
}
