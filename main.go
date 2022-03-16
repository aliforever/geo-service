package main

import (
	"fmt"
	"github.com/aliforever/geo-service/geolocation"
	"time"
)

func main() {
	n := time.Now()

	// go func() {
	// 	for bytes := range receiver {
	//
	// 		// fmt.Printf("%+v\n", info)
	// 	}
	// }()
	path := "data_dump.csv"
	_, err := ParseCSV(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	then := time.Now()
	fmt.Printf("Took %s to complete\n", then.Sub(n).String())
}

func parseData(data []byte) {
	_, err := geolocation.NewGeoLocationFromRowBytes(data)
	if err != nil {
		// if err != emptyIPAddress && err != emptyLat && err != emptyLong && !strings.Contains(err.Error(), "invalid syntax") {
		// 	fmt.Println(err)
		// 	panic(string(bytes))
		// }
		return
	}
}
