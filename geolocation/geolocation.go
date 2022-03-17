package geolocation

import (
	"net"
)

// GeoLocation ip_address,country_code,country,city,latitude,longitude,mystery_value
type GeoLocation struct {
	IPAddress    net.IP  `json:"ip_address"`
	CountryCode  string  `json:"country_code"`
	Country      string  `json:"country"`
	City         string  `json:"city"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	MysteryValue int64   `json:"mystery_value"`
}

func NewGeoLocationFromString(data string) (g *GeoLocation, err error) {
	var (
		ipAddr                     net.IP
		countryCode, country, city string
		lat, lng                   float64
		mysteryValue               int64
	)

	ipAddr, countryCode, country, city, lat, lng, mysteryValue, err = parseColumns(getColumns(data))
	if err != nil {
		return
	}

	g = &GeoLocation{
		IPAddress:    ipAddr,
		CountryCode:  countryCode,
		Country:      country,
		City:         city,
		Latitude:     lat,
		Longitude:    lng,
		MysteryValue: mysteryValue,
	}

	return
}
