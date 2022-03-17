package geoservice

import (
	"bufio"
	"github.com/aliforever/geo-service/geolocation"
	"net"
	"os"
	"sync"
	"time"
)

type GeoService struct {
	db geolocation.Repository
}

func NewGeoService(db geolocation.Repository) (gs *GeoService) {
	return &GeoService{db: db}
}

// initializeWorker receives number of workers defining number of goroutines for initializing GeoLocation from rows
// And writing the results to the ch channel
func (g *GeoService) initializeWorker(workers int, rows [][]byte, ch chan *geolocation.GeoLocation) {
	var wg sync.WaitGroup

	goroutines := len(rows) / workers

	var firstIndex, lastIndex int
	for i := 0; i < goroutines; i++ {

		firstIndex = i * workers
		lastIndex = firstIndex + workers

		values := rows[firstIndex:lastIndex]
		if i == goroutines-1 {
			values = rows[firstIndex:]
		}

		wg.Add(1)
		go func(rows [][]byte) {
			defer wg.Done()
			for _, row := range rows {
				loc, locErr := geolocation.NewGeoLocationFromRowBytes(row)
				if locErr != nil || loc == nil {
					continue
				}
				ch <- loc
			}
		}(values)
	}

	wg.Wait()
	close(ch)
}

func (g *GeoService) ParseCSV(path string, workers int) (locations []*geolocation.GeoLocation, stat *Statistics, err error) {
	begin := time.Now()

	var file *os.File
	file, err = os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return
	}
	defer file.Close()

	var fileStat os.FileInfo
	fileStat, err = file.Stat()
	if err != nil {
		return
	}

	var maxSize int
	maxSize = int(fileStat.Size())
	buffer := make([]byte, 0, maxSize)

	r := bufio.NewScanner(file)
	r.Buffer(buffer, maxSize)

	r.Scan() // This is to skip header row

	var rowChan = make(chan *geolocation.GeoLocation)

	var rows [][]byte
	for r.Scan() {
		rows = append(rows, r.Bytes())
	}

	err = r.Err()
	if err != nil {
		return
	}

	appendBegin := time.Now()
	var appendElapsed time.Duration

	var duplicates int

	var wg sync.WaitGroup
	go func() {
		wg.Add(1)
		defer wg.Done()

		var storage = map[string]*geolocation.GeoLocation{}
		for location := range rowChan {
			if storage[location.IPAddress.String()] != nil {
				duplicates++
				continue
			}
			storage[location.IPAddress.String()] = location
		}

		for _, location := range storage {
			locations = append(locations, location)
		}

		appendElapsed = time.Now().Sub(appendBegin)
	}()

	parsedBegin := time.Now()
	g.initializeWorker(workers, rows, rowChan)
	parsedElapsed := time.Now().Sub(parsedBegin)

	wg.Wait()

	end := time.Now()

	stat = &Statistics{
		Elapsed:          end.Sub(begin),
		ElapsedParsed:    parsedElapsed,
		ElapsedAppend:    appendElapsed,
		Duplicates:       duplicates,
		AcceptedEntries:  len(locations),
		DiscardedEntries: len(rows) - len(locations),
	}
	return
}

func (g *GeoService) StoreLocations(locations []*geolocation.GeoLocation) (err error) {
	for _, location := range locations {
		err = g.db.Store(location)
		if err != nil {
			return
		}
	}
	return
}

func (g *GeoService) RetrieveLocation(ip net.IP) (location *geolocation.GeoLocation, err error) {
	location, err = g.db.Retrieve(ip)
	return
}
