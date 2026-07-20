package zippopotam

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"discordgo-skeleton/internal/application/ports"
	"discordgo-skeleton/internal/domain/location"

	"discordgo-skeleton/pkg/httpx"
)

var _ ports.Geocoder = (*geocoder)(nil)

const lookupURL = "https://api.zippopotam.us/us"

type lookupResponse struct {
	PostCode string `json:"post code"`
	Places   []struct {
		PlaceName string `json:"place name"`
		State     string `json:"state"`
		Latitude  string `json:"latitude"`
		Longitude string `json:"longitude"`
	} `json:"places"`
}

type geocoder struct {
	client *httpx.Client
}

func NewGeocoder() *geocoder {
	return &geocoder{
		client: httpx.NewClient(),
	}
}

func (g *geocoder) Lookup(zip string) (location.Location, error) {
	resp, err := g.client.Get(fmt.Sprintf("%s/%s", lookupURL, zip))
	if err != nil {
		return location.Location{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return location.Location{}, location.ErrLocationNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return location.Location{}, fmt.Errorf("geocoder returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return location.Location{}, err
	}

	var data lookupResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return location.Location{}, err
	}
	if len(data.Places) == 0 {
		return location.Location{}, location.ErrLocationNotFound
	}

	place := data.Places[0]
	lat, err := strconv.ParseFloat(place.Latitude, 64)
	if err != nil {
		return location.Location{}, fmt.Errorf("parse latitude: %w", err)
	}
	lon, err := strconv.ParseFloat(place.Longitude, 64)
	if err != nil {
		return location.Location{}, fmt.Errorf("parse longitude: %w", err)
	}

	return location.New(data.PostCode, place.PlaceName, place.State, lat, lon), nil
}
