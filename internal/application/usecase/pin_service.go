package usecase

import (
	"errors"
	"fmt"

	"discordgo-skeleton/internal/application/ports"
	"discordgo-skeleton/internal/domain/location"
	"discordgo-skeleton/internal/domain/pin"
)

var (
	ErrAddZipInvalid  = errors.New("invalid zip code")
	ErrAddDuplicate   = errors.New("zip code already pinned")
	ErrRemoveNotFound = errors.New("zip code not pinned")
)

type PinService struct {
	repository pin.Repository
	geocoder   ports.Geocoder
}

func NewPinService(repository pin.Repository, geocoder ports.Geocoder) *PinService {
	return &PinService{
		repository: repository,
		geocoder:   geocoder,
	}
}

func (s *PinService) Add(userID string, zip string) (location.Location, error) {
	loc, err := s.geocoder.Lookup(zip)
	if err != nil {
		if errors.Is(err, location.ErrLocationNotFound) {
			return location.Location{}, fmt.Errorf("%w: %s", ErrAddZipInvalid, zip)
		}
		return location.Location{}, err
	}

	newPin := pin.New(userID, zip)
	if err := s.repository.Create(newPin); err != nil {
		if errors.Is(err, pin.ErrPinDuplicate) {
			return location.Location{}, fmt.Errorf("%w: %s", ErrAddDuplicate, zip)
		}
		return location.Location{}, err
	}
	return loc, nil
}

func (s *PinService) Remove(userID string, zip string) (location.Location, error) {
	if err := s.repository.Delete(userID, zip); err != nil {
		if errors.Is(err, pin.ErrPinNotFound) {
			return location.Location{}, fmt.Errorf("%w: %s", ErrRemoveNotFound, zip)
		}
		return location.Location{}, err
	}

	loc, err := s.geocoder.Lookup(zip)
	if err != nil {
		return location.Location{}, err
	}

	return loc, nil
}
