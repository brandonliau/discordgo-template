package ports

import (
	"discordgo-skeleton/internal/domain/location"
	"discordgo-skeleton/internal/domain/weather"
)

type WeatherFeed interface {
	Fetch(loc location.Location) (weather.Weather, error)
}
