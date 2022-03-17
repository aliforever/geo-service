package geoservice

import (
	"errors"
	"github.com/aliforever/geo-service/geolocation"
	"net"
	"sync"
)

// testDB This is a custom db to test package repository (Store, Retrieve Methods)
type testDB struct {
	sync.Mutex
	data map[string]*geolocation.GeoLocation
}

func newTestDB() *testDB {
	return &testDB{
		Mutex: sync.Mutex{},
		data:  map[string]*geolocation.GeoLocation{},
	}
}

func (t *testDB) Store(g *geolocation.GeoLocation) (err error) {
	t.Lock()
	defer t.Unlock()

	if _, ok := t.data[g.IPAddress.String()]; ok {
		err = errors.New("data exists")
		return
	}

	t.data[g.IPAddress.String()] = g
	return
}

func (t *testDB) StoreMany(gs []*geolocation.GeoLocation) (err error) {
	t.Lock()
	defer t.Unlock()

	for _, g := range gs {
		if _, ok := t.data[g.IPAddress.String()]; ok {
			err = errors.New("data exists")
			return
		}
		t.data[g.IPAddress.String()] = g
	}

	return
}

func (t *testDB) Retrieve(ip net.IP) (g *geolocation.GeoLocation, err error) {
	t.Lock()
	defer t.Unlock()

	if g = t.data[ip.String()]; g == nil {
		err = errors.New("data does not exists")
		return
	}

	return
}
