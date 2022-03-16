package geolocation

import (
	"errors"
	"net"
	"regexp"
	"strconv"
	"strings"
)

var (
	invalidDataError = errors.New("invalid_data")
	emptyIPAddress   = errors.New("empty_ip_address")
	invalidIPAddress = errors.New("invalid_ip_address")
	emptyLat         = errors.New("empty_latitude")
	emptyLong        = errors.New("empty_longitude")
)

func escapeCommas(data string) string {
	r := regexp.MustCompile(`(".*?,.*?")|('.*?,.*?')`)

	for _, match := range r.FindAllString(data, -1) {
		escaped := strings.ReplaceAll(match, ",", "-")
		data = strings.ReplaceAll(data, match, escaped)
	}

	return data
}

func getColumns(data []byte) (columns []string) {
	trimmed := strings.TrimSpace(string(data))
	columns = strings.Split(trimmed, ",")
	for index := range columns {
		columns[index] = strings.TrimSpace(columns[index])
	}

	// If we have more than 7 columns then it probably means there are unescaped commas, inside Country or City columns.
	// We try to escape commas and re-calculate column length
	if len(columns) > 7 {
		columns = strings.Split(escapeCommas(trimmed), ",")
	}

	return
}

// parseColumns
// ==== Rules ====
// The IPAddress can't be Empty
// Column Length Should Match 7
// Latitude & Longitude can't be Empty
// Latitude & Longitude should be of type Float
// Mystery Value should be of type Int
func parseColumns(columns []string) (ipAddr net.IP, countryCode, country, city string, lat, lng float64, mysteryValue int64, err error) {
	if columns[0] == "" {
		err = emptyIPAddress
		return
	}

	if ipAddr = net.ParseIP(columns[0]); ipAddr == nil {
		err = invalidIPAddress
		return
	}

	if len(columns) != 7 {
		err = invalidDataError
		return
	}

	if columns[4] == "" {
		err = emptyLat
		return
	}

	if columns[5] == "" {
		err = emptyLong
		return
	}

	lat, err = strconv.ParseFloat(columns[4], 64)
	if err != nil {
		return
	}

	lng, err = strconv.ParseFloat(columns[5], 64)
	if err != nil {
		return
	}

	mysteryValue, err = strconv.ParseInt(columns[6], 0, 64)
	if err != nil {
		return
	}

	countryCode = columns[1]
	country = columns[2]
	city = columns[3]

	return
}
