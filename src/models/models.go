package models

// Stores the geo coordinates in lat, lon format.
type Coordinate struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

// Response received for conversion of city name to geo co-ordinates.
type CityQueryResponseModel struct {
	Name        string
	Local_names interface{}
	Lat         float64
	Lon         float64
	Country     string
	State       string
}
