package weather

type Weather struct {
	Code       WeatherCode
	TempF      float64
	FeelsLikeF float64
	Humidity   int
	WindMph    float64
	WindDeg    float64
	PrecipIn   float64
	IsDay      bool
	HighF      float64
	LowF       float64
}

type WeatherCode int

func (c WeatherCode) Description() string {
	switch c {
	case 0:
		return "Clear sky"
	case 1:
		return "Mainly clear"
	case 2:
		return "Partly cloudy"
	case 3:
		return "Overcast"
	case 45:
		return "Fog"
	case 48:
		return "Depositing rime fog"
	case 51:
		return "Light drizzle"
	case 53:
		return "Moderate drizzle"
	case 55:
		return "Heavy drizzle"
	case 56:
		return "Light freezing drizzle"
	case 57:
		return "Heavy freezing drizzle"
	case 61:
		return "Slight rain"
	case 63:
		return "Moderate rain"
	case 65:
		return "Heavy rain"
	case 71:
		return "Slight snow fall"
	case 73:
		return "Moderate snow fall"
	case 75:
		return "Heavy snow fall"
	case 77:
		return "Snow grains"
	case 80:
		return "Slight rain showers"
	case 81:
		return "Moderate rain showers"
	case 82:
		return "Violent rain showers"
	case 85:
		return "Slight snow showers"
	case 86:
		return "Heavy snow showers"
	case 95:
		return "Thunderstorm"
	case 96:
		return "Thunderstorm with slight hail"
	case 99:
		return "Thunderstorm with heavy hail"
	default:
		return "Unknown"
	}
}
