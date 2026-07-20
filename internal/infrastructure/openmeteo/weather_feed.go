package openmeteo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"discordgo-skeleton/internal/application/ports"
	"discordgo-skeleton/internal/domain/location"
	"discordgo-skeleton/internal/domain/weather"

	"discordgo-skeleton/pkg/httpx"
)

var _ ports.WeatherFeed = (*weatherFeed)(nil)

const forecastURL = "https://api.open-meteo.com/v1/forecast"

type forecastResponse struct {
	Current struct {
		Temperature   float64 `json:"temperature_2m"`
		ApparentTemp  float64 `json:"apparent_temperature"`
		Humidity      int     `json:"relative_humidity_2m"`
		WindSpeed     float64 `json:"wind_speed_10m"`
		WindDirection float64 `json:"wind_direction_10m"`
		Precipitation float64 `json:"precipitation"`
		WeatherCode   int     `json:"weather_code"`
		IsDay         int     `json:"is_day"`
	} `json:"current"`
	Daily struct {
		TempMax []float64 `json:"temperature_2m_max"`
		TempMin []float64 `json:"temperature_2m_min"`
	} `json:"daily"`
}

type weatherFeed struct {
	client *httpx.Client
}

func NewWeatherFeed() *weatherFeed {
	return &weatherFeed{
		client: httpx.NewClient(),
	}
}

func (f *weatherFeed) Fetch(loc location.Location) (weather.Weather, error) {
	params := url.Values{}
	params.Set("latitude", strconv.FormatFloat(loc.Latitude, 'f', -1, 64))
	params.Set("longitude", strconv.FormatFloat(loc.Longitude, 'f', -1, 64))
	params.Set("current", "temperature_2m,relative_humidity_2m,apparent_temperature,is_day,precipitation,weather_code,wind_speed_10m,wind_direction_10m")
	params.Set("daily", "temperature_2m_max,temperature_2m_min")
	params.Set("temperature_unit", "fahrenheit")
	params.Set("wind_speed_unit", "mph")
	params.Set("precipitation_unit", "inch")
	params.Set("timezone", "auto")
	params.Set("forecast_days", "1")

	resp, err := f.client.Get(fmt.Sprintf("%s?%s", forecastURL, params.Encode()))
	if err != nil {
		return weather.Weather{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return weather.Weather{}, fmt.Errorf("weather feed returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return weather.Weather{}, err
	}

	var data forecastResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return weather.Weather{}, err
	}

	w := weather.Weather{
		Code:       weather.WeatherCode(data.Current.WeatherCode),
		TempF:      data.Current.Temperature,
		FeelsLikeF: data.Current.ApparentTemp,
		Humidity:   data.Current.Humidity,
		WindMph:    data.Current.WindSpeed,
		WindDeg:    data.Current.WindDirection,
		PrecipIn:   data.Current.Precipitation,
		IsDay:      data.Current.IsDay == 1,
	}
	if len(data.Daily.TempMax) > 0 {
		w.HighF = data.Daily.TempMax[0]
	}
	if len(data.Daily.TempMin) > 0 {
		w.LowF = data.Daily.TempMin[0]
	}

	return w, nil
}
