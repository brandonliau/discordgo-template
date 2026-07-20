package view

import (
	"fmt"

	"discordgo-skeleton/internal/domain/location"
	"discordgo-skeleton/internal/domain/weather"
)

type WeatherView struct {
	Zip        string
	City       string
	State      string
	Condition  string
	IsDay      bool
	TempF      float64
	FeelsLikeF float64
	HighF      float64
	LowF       float64
	Humidity   int
	WindMph    float64
	WindDir    string
}

func FromWeather(loc location.Location, w weather.Weather) WeatherView {
	return WeatherView{
		Zip:        loc.Zip,
		City:       loc.City,
		State:      loc.State,
		Condition:  w.Code.Description(),
		IsDay:      w.IsDay,
		TempF:      w.TempF,
		FeelsLikeF: w.FeelsLikeF,
		HighF:      w.HighF,
		LowF:       w.LowF,
		Humidity:   w.Humidity,
		WindMph:    w.WindMph,
		WindDir:    compassDirection(w.WindDeg),
	}
}

func (v WeatherView) Location() string {
	place := v.City
	if v.State != "" {
		place = fmt.Sprintf("%s, %s", v.City, v.State)
	}
	if v.Zip != "" {
		return fmt.Sprintf("%s (%s)", place, v.Zip)
	}
	return place
}

func compassDirection(deg float64) string {
	points := []string{"N", "NE", "E", "SE", "S", "SW", "W", "NW"}
	idx := int((deg+22.5)/45.0) % len(points)
	if idx < 0 {
		idx += len(points)
	}
	return points[idx]
}
