package geolocation

import "net"

type Repository interface {
	Store(*GeoLocation) error
	StoreMany([]*GeoLocation) error
	Retrieve(ipAddress net.IP) (*GeoLocation, error)
}
