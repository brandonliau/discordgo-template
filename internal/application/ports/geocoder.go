package ports

import (
	"discordgo-skeleton/internal/domain/location"
)

type Geocoder interface {
	Lookup(zip string) (location.Location, error)
}
