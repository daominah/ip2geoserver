package main

import (
	"github.com/daominah/echo_ip_httpsvr/ip2geo"
	"github.com/mywrap/log"
)

func main() {
	geoReader, err := ip2geo.NewDefaultReader()
	if err != nil {
		log.Fatal(err)
	}
	geoInfo, err := geoReader.ReadIPInfo("1.1.1.1")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("geoInfo: %+v", geoInfo)
}
