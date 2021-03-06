package geoservice

import (
	"github.com/aliforever/geo-service/geolocation"
	"io/ioutil"
	"net"
	"os"
	"reflect"
	"testing"
)

func compareLocations(locs1, locs2 []*geolocation.GeoLocation) bool {
	for _, location := range locs1 {
		isFound := false
		for _, geoLocation := range locs2 {
			if geoLocation.IPAddress.String() == location.IPAddress.String() {
				isFound = true
				if location.Longitude != geoLocation.Longitude || location.Latitude != geoLocation.Latitude ||
					location.City != geoLocation.City || location.Country != geoLocation.Country ||
					location.CountryCode != geoLocation.CountryCode {
					return false
				}
			}
		}
		if !isFound {
			return false
		}
	}

	return true
}

func BenchmarkGeoService_ParseCSV(b *testing.B) {
	info := "ip_address,country_code,country,city,latitude,longitude,mystery_value\n"
	info += "201.106.141.15,SI,Nepal,DuBuquemouth,-84.87503094689836,7.206435933364332,7823011346\n"
	info += "70.95.73.73,TL,Saudi Arabia,Gradymouth,-49.16675918861615,-86.05920084416894,2559997162\n"
	info += "160.103.7.140,CZ,Nicaragua,New Neva,-68.31023296602508,-37.62435199624531,7301823115\n"
	info += ",PY,Falkland Islands (Malvinas),,75.41685191518815,-144.6943217219469,0"
	err := ioutil.WriteFile("data_dump1.csv", []byte(info), 0644)
	if err != nil {
		b.Errorf("ParseCSV() error = cant write test data: %s", err)
		return
	}
	defer os.Remove("data_dump1.csv")

	for n := 0; n < b.N; n++ {
		g := &GeoService{}
		g.ParseCSV("data_dump1.csv", 5)
	}

}

func TestGeoService_ParseCSV(t *testing.T) {
	info := "ip_address,country_code,country,city,latitude,longitude,mystery_value\n"
	info += "201.106.141.15,SI,Nepal,DuBuquemouth,-84.87503094689836,7.206435933364332,7823011346\n"
	info += "70.95.73.73,TL,Saudi Arabia,Gradymouth,-49.16675918861615,-86.05920084416894,2559997162\n"
	info += "160.103.7.140,CZ,Nicaragua,New Neva,-68.31023296602508,-37.62435199624531,7301823115\n"
	info += ",PY,Falkland Islands (Malvinas),,75.41685191518815,-144.6943217219469,0"
	err := ioutil.WriteFile("data_dump1.csv", []byte(info), 0644)
	if err != nil {
		t.Errorf("ParseCSV() error = cant write test data: %s", err)
		return
	}
	defer os.Remove("data_dump1.csv")

	type args struct {
		path    string
		workers int
	}

	tests := []struct {
		name          string
		args          args
		wantLocations []*geolocation.GeoLocation
		wantStat      *Statistics
		wantErr       bool
	}{
		{
			name: "Test1",
			args: args{path: "data_dump1.csv", workers: 5},
			wantLocations: []*geolocation.GeoLocation{
				{
					IPAddress:    net.ParseIP("201.106.141.15"),
					CountryCode:  "SI",
					Country:      "Nepal",
					City:         "DuBuquemouth",
					Latitude:     -84.87503094689836,
					Longitude:    7.206435933364332,
					MysteryValue: 7823011346,
				},
				{
					IPAddress:    net.ParseIP("160.103.7.140"),
					CountryCode:  "CZ",
					Country:      "Nicaragua",
					City:         "New Neva",
					Latitude:     -68.31023296602508,
					Longitude:    -37.62435199624531,
					MysteryValue: 7301823115,
				},
				{
					IPAddress:    net.ParseIP("70.95.73.73"),
					CountryCode:  "TL",
					Country:      "Saudi Arabia",
					City:         "Gradymouth",
					Latitude:     -49.16675918861615,
					Longitude:    -86.05920084416894,
					MysteryValue: 2559997162,
				},
			},
			wantStat: &Statistics{
				AcceptedEntries:  3,
				DiscardedEntries: 1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GeoService{}
			gotLocations, gotStat, err := g.ParseCSV(tt.args.path, tt.args.workers)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCSV() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !compareLocations(gotLocations, tt.wantLocations) {
				t.Errorf("ParseCSV() gotLocations = %v, want %v", gotLocations, tt.wantLocations)
			}
			if (gotStat.AcceptedEntries != tt.wantStat.AcceptedEntries) || (gotStat.DiscardedEntries != tt.wantStat.DiscardedEntries) {
				t.Errorf("ParseCSV() gotStatsAcceptedEnteries = %v, want %v", gotStat, tt.wantStat)
			}
		})
	}

}

func TestGeoService_StoreLocations(t *testing.T) {
	err := ioutil.WriteFile("data_dump1.csv", []byte(`ip_address,country_code,country,city,latitude,longitude,mystery_value
200.106.141.15,SI,Nepal,DuBuquemouth,-84.87503094689836,7.206435933364332,7823011346
160.103.7.140,CZ,Nicaragua,New Neva,-68.31023296602508,-37.62435199624531,7301823115
70.95.73.73,TL,Saudi Arabia,Gradymouth,-49.16675918861615,-86.05920084416894,2559997162
,PY,Falkland Islands (Malvinas),,75.41685191518815,-144.6943217219469,0`), 0644)
	if err != nil {
		t.Errorf("ParseCSV() error = cant write test data: %s", err)
		return
	}
	defer os.Remove("data_dump1.csv")

	db := newTestDB()

	type args struct {
		db        geolocation.Repository
		locations []*geolocation.GeoLocation
	}
	tests := []struct {
		name      string
		args      args
		wantStats *Statistics
		wantErr   bool
	}{
		{
			name: "Test1",
			args: args{
				db: db,
				locations: []*geolocation.GeoLocation{
					{
						IPAddress:    net.ParseIP("200.106.141.15"),
						CountryCode:  "SI",
						Country:      "Nepal",
						City:         "DuBuquemouth",
						Latitude:     -84.87503094689836,
						Longitude:    7.206435933364332,
						MysteryValue: 7823011346,
					},
					{
						IPAddress:    net.ParseIP("160.103.7.140"),
						CountryCode:  "CZ",
						Country:      "Nicaragua",
						City:         "New Neva",
						Latitude:     -68.31023296602508,
						Longitude:    -37.62435199624531,
						MysteryValue: 7301823115,
					},
					{
						IPAddress:    net.ParseIP("70.95.73.73"),
						CountryCode:  "TL",
						Country:      "Saudi Arabia",
						City:         "Gradymouth",
						Latitude:     -49.16675918861615,
						Longitude:    -86.05920084416894,
						MysteryValue: 2559997162,
					},
				},
			},
			wantStats: &Statistics{
				AcceptedEntries:  3,
				DiscardedEntries: 0,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGeoService(tt.args.db)
			err := g.StoreLocations(tt.args.locations)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreLocations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGeoService_RetrieveLocation(t *testing.T) {
	err := ioutil.WriteFile("data_dump1.csv", []byte(`ip_address,country_code,country,city,latitude,longitude,mystery_value
200.106.141.15,SI,Nepal,DuBuquemouth,-84.87503094689836,7.206435933364332,7823011346
160.103.7.140,CZ,Nicaragua,New Neva,-68.31023296602508,-37.62435199624531,7301823115
70.95.73.73,TL,Saudi Arabia,Gradymouth,-49.16675918861615,-86.05920084416894,2559997162
,PY,Falkland Islands (Malvinas),,75.41685191518815,-144.6943217219469,0`), 0644)
	if err != nil {
		t.Errorf("ParseCSV() error = cant write test data: %s", err)
		return
	}
	defer os.Remove("data_dump1.csv")

	db := newTestDB()

	g := NewGeoService(db)
	locations, _, err := g.ParseCSV("data_dump1.csv", 5)
	err = g.StoreLocations(locations)
	if err != nil {
		t.Errorf("cant store locations: %s", err)
		return
	}

	type args struct {
		db geolocation.Repository
		ip net.IP
	}
	tests := []struct {
		name         string
		args         args
		wantLocation *geolocation.GeoLocation
		wantErr      bool
	}{
		{
			name: "Test1",
			args: args{
				db: db,
				ip: net.ParseIP("200.106.141.15"),
			},
			wantLocation: &geolocation.GeoLocation{
				IPAddress:    net.ParseIP("200.106.141.15"),
				CountryCode:  "SI",
				Country:      "Nepal",
				City:         "DuBuquemouth",
				Latitude:     -84.87503094689836,
				Longitude:    7.206435933364332,
				MysteryValue: 7823011346,
			},
			wantErr: false,
		},
		{
			name: "Test2",
			args: args{
				db: db,
				ip: net.ParseIP("200.20.141.16"),
			},
			wantLocation: nil,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGeoService(tt.args.db)
			gotLocation, err := g.RetrieveLocation(tt.args.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("RetrieveLocation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotLocation, tt.wantLocation) {
				t.Errorf("RetrieveLocation() gotLocation = %v, want %v", gotLocation, tt.wantLocation)
			}
		})
	}
}
