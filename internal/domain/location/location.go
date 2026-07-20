package location

type Location struct {
	Zip       string
	City      string
	State     string
	Latitude  float64
	Longitude float64
}

func New(zip string, city string, state string, latitude float64, longitude float64) Location {
	return Location{
		Zip:       zip,
		City:      city,
		State:     state,
		Latitude:  latitude,
		Longitude: longitude,
	}
}
