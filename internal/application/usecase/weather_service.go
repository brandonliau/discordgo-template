package usecase

import (
	"errors"
	"fmt"

	"discordgo-skeleton/internal/application/ports"
	"discordgo-skeleton/internal/application/view"
	"discordgo-skeleton/internal/domain/location"
	"discordgo-skeleton/internal/domain/pin"
)

var (
	ErrSearchZipInvalid = errors.New("invalid zip code")
)

type WeatherService struct {
	geocoder    ports.Geocoder
	weatherFeed ports.WeatherFeed
	pins        pin.Repository
}

func NewWeatherService(geocoder ports.Geocoder, weatherFeed ports.WeatherFeed, pins pin.Repository) *WeatherService {
	return &WeatherService{
		geocoder:    geocoder,
		weatherFeed: weatherFeed,
		pins:        pins,
	}
}

func (s *WeatherService) Search(zip string) (view.WeatherView, error) {
	loc, err := s.geocoder.Lookup(zip)
	if err != nil {
		if errors.Is(err, location.ErrLocationNotFound) {
			return view.WeatherView{}, fmt.Errorf("%w: %s", ErrSearchZipInvalid, zip)
		}
		return view.WeatherView{}, err
	}

	w, err := s.weatherFeed.Fetch(loc)
	if err != nil {
		return view.WeatherView{}, err
	}
	return view.FromWeather(loc, w), nil
}

func (s *WeatherService) Random() (view.WeatherView, error) {
	loc := location.RandomLocation()
	w, err := s.weatherFeed.Fetch(loc)
	if err != nil {
		return view.WeatherView{}, err
	}
	return view.FromWeather(loc, w), nil
}

func (s *WeatherService) List(userID string) ([]view.WeatherView, error) {
	pins, err := s.pins.ListByUser(userID)
	if err != nil {
		return nil, err
	}

	views := make([]view.WeatherView, 0, len(pins))
	for _, pinned := range pins {
		resolved, err := s.geocoder.Lookup(pinned.Zip)
		if err != nil {
			return nil, err
		}
		w, err := s.weatherFeed.Fetch(resolved)
		if err != nil {
			return nil, err
		}
		views = append(views, view.FromWeather(resolved, w))
	}
	return views, nil
}
