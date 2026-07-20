package location

import (
	_ "embed"
	"encoding/json"
	"math/rand"
)

//go:embed locations.json
var locationsData []byte

var randomLocations = mustLoadRandomLocations()

func mustLoadRandomLocations() []Location {
	var locations []Location
	if err := json.Unmarshal(locationsData, &locations); err != nil {
		panic(err)
	}
	return locations
}

func RandomLocation() Location {
	return randomLocations[rand.Intn(len(randomLocations))]
}
