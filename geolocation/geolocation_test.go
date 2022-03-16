package geolocation

import (
	"testing"
)

func Test_escapeCommas(t *testing.T) {
	type args struct {
		data string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test1",
			args: args{data: `152.159.31.208,GA,"Virgin Islands, British",'Lake, Wavatown',12.964804277773922,-56.656208830174734,1878158074`},
			want: `152.159.31.208,GA,"Virgin Islands- British",'Lake- Wavatown',12.964804277773922,-56.656208830174734,1878158074`,
		},
		{
			name: "Test2",
			args: args{data: `152.159.31.208,GA,"Virgin Islands, British",Lake Wavatown,12.964804277773922,-56.656208830174734,1878158074`},
			want: `152.159.31.208,GA,"Virgin Islands- British",Lake Wavatown,12.964804277773922,-56.656208830174734,1878158074`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := escapeCommas(tt.args.data); got != tt.want {
				t.Errorf("escapeCommas()\nGot: %v\nWant: %v", got, tt.want)
			}
		})
	}
}

func Test_parseColumns(t *testing.T) {
	type args struct {
		columns []string
	}
	tests := []struct {
		name             string
		args             args
		wantIpAddr       string
		wantCountryCode  string
		wantCountry      string
		wantCity         string
		wantLat          float64
		wantLng          float64
		wantMysteryValue int64
		wantErr          bool
	}{
		{
			name:             "Test1",
			args:             args{columns: getColumns([]byte(`200.106.141.15,SI,Nepal,DuBuquemouth,-84.87503094689836,7.206435933364332,7823011346`))},
			wantIpAddr:       "200.106.141.15",
			wantCountryCode:  "SI",
			wantCountry:      "Nepal",
			wantCity:         "DuBuquemouth",
			wantLat:          -84.87503094689836,
			wantLng:          7.206435933364332,
			wantMysteryValue: 7823011346,
			wantErr:          false,
		},
		{
			name:             "Test2",
			args:             args{columns: getColumns([]byte(`160.103.7.140,CZ,Nicaragua,New Neva,-68.31023296602508,-37.62435199624531,7301823115`))},
			wantIpAddr:       "160.103.7.140",
			wantCountryCode:  "CZ",
			wantCountry:      "Nicaragua",
			wantCity:         "New Neva",
			wantLat:          -68.31023296602508,
			wantLng:          -37.62435199624531,
			wantMysteryValue: 7301823115,
			wantErr:          false,
		},
		{
			name:             "Test3",
			args:             args{columns: getColumns([]byte(`160.103.7.140,CZ,Nicaragua,New Neva,-68.31023296602508,-37.62435199624531,7301823115,`))},
			wantIpAddr:       "160.103.7.140",
			wantCountryCode:  "CZ",
			wantCountry:      "Nicaragua",
			wantCity:         "New Neva",
			wantLat:          -68.31023296602508,
			wantLng:          -37.62435199624531,
			wantMysteryValue: 7301823115,
			wantErr:          true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIpAddr, gotCountryCode, gotCountry, gotCity, gotLat, gotLng, gotMysteryValue, err := parseColumns(tt.args.columns)
			if (err != nil) && !tt.wantErr {
				t.Errorf("parseColumns() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if gotIpAddr != tt.wantIpAddr {
					t.Errorf("parseColumns() gotIpAddr = %v, want %v", gotIpAddr, tt.wantIpAddr)
				}
				if gotCountryCode != tt.wantCountryCode {
					t.Errorf("parseColumns() gotCountryCode = %v, want %v", gotCountryCode, tt.wantCountryCode)
				}
				if gotCountry != tt.wantCountry {
					t.Errorf("parseColumns() gotCountry = %v, want %v", gotCountry, tt.wantCountry)
				}
				if gotCity != tt.wantCity {
					t.Errorf("parseColumns() gotCity = %v, want %v", gotCity, tt.wantCity)
				}
				if gotLat != tt.wantLat {
					t.Errorf("parseColumns() gotLat = %v, want %v", gotLat, tt.wantLat)
				}
				if gotLng != tt.wantLng {
					t.Errorf("parseColumns() gotLng = %v, want %v", gotLng, tt.wantLng)
				}
				if gotMysteryValue != tt.wantMysteryValue {
					t.Errorf("parseColumns() gotMysteryValue = %v, want %v", gotMysteryValue, tt.wantMysteryValue)
				}
			}
		})
	}
}
