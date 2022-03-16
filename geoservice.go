package main

import (
	"bufio"
	"github.com/aliforever/geo-service/geolocation"
	"os"
	"time"
)

type GeoService struct {
	db geolocation.Repository
}

func NewGeoService() (gs *GeoService) {
	return &GeoService{}
}

func (g *GeoService) ParseCSV(path string) (locations []*geolocation.GeoLocation, stat *Statistics, err error) {
	begin := time.Now()

	var file *os.File
	file, err = os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return
	}
	defer file.Close()

	r := bufio.NewScanner(file)
	r.Scan() // This is to skip header row

	var rows [][]byte
	// var data []byte
	for r.Scan() {
		rows = append(rows, r.Bytes())
	}

	err = r.Err()
	if err != nil {
		return
	}

	for _, row := range rows {
		loc, locErr := geolocation.NewGeoLocationFromRowBytes(row)
		if locErr != nil {
			continue
		}
		locations = append(locations, loc)
	}

	end := time.Now()

	stat = &Statistics{
		TimeElapsed:      begin.Sub(end),
		AcceptedEntries:  len(locations),
		DiscardedEntries: len(rows) - len(locations),
	}
	return
}

func (g *GeoService) StoreLocations(db geolocation.Repository, locations []*geolocation.GeoLocation) (stats *Statistics, err error) {
	begin := time.Now()

	var accepted, discarded int
	for _, location := range locations {
		storeErr := db.Store(location)
		if storeErr != nil {
			discarded++
			continue
		}
		accepted++
	}

	end := time.Now()

	stats = &Statistics{
		TimeElapsed:      end.Sub(begin),
		AcceptedEntries:  accepted,
		DiscardedEntries: discarded,
	}
	return
}
