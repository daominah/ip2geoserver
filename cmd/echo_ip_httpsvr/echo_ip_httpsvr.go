package main

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"

	"github.com/daominah/echo_ip_httpsvr/ip2geo"
	"github.com/mywrap/httpsvr"
	"github.com/mywrap/log"
)

func handlerEcho(w http.ResponseWriter, r *http.Request) {
	log.Printf("____________________________________________________\n")
	rDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		w.WriteHeader(500)
		httpsvr.Write(w, r, fmt.Sprintf("error DumpRequest: %v", err))
		return
	}
	httpsvr.Write(w, r, fmt.Sprintf("echo req from %v:\n\n", r.RemoteAddr))
	httpsvr.Write(w, r, string(rDump))
	log.Printf("____________________________________________________\n")
	return
}

func handlerRawIP(w http.ResponseWriter, r *http.Request) {
	ip := ip2geo.GetIpFromAddress(r.RemoteAddr)
	httpsvr.Write(w, r, fmt.Sprintf("%v", ip))
	return
}

func handlerGeoIP() http.HandlerFunc {
	geoReader, err := ip2geo.NewDefaultReader()
	if err != nil {
		log.Fatal(err)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ip0 := httpsvr.GetUrlParams(r)["ip0"]
		if ip0 == "" {
			ip0 = ip2geo.GetIpFromAddress(r.RemoteAddr) // host:port -> host
		} else { // caller can input an IP or hostname
			if net.ParseIP(ip0) == nil {
				ip0 = ip2geo.LookupIPFromHost(ip0)
			}
		}
		geo, err := geoReader.ReadIPInfo(ip0)
		if err != nil {
			w.WriteHeader(400)
			httpsvr.Write(w, r, fmt.Sprintf(
				"error ReadIPInfo %v: %v", ip0, err))
			return
		}
		httpsvr.WriteJson(w, r, geo)
		return
	}
}

func main() {
	server := httpsvr.NewServer()
	initedHandlerGeoIP := handlerGeoIP()
	server.AddHandler("GET", "/ip", handlerRawIP)
	server.AddHandler("GET", "/ip/geo", initedHandlerGeoIP)
	server.AddHandler("GET", "/ip/geo/:ip0", initedHandlerGeoIP)
	server.AddHandlerNotFound(handlerEcho)

	listen := os.Getenv("LISTENING_PORT")
	if listen == "" {
		listen = ":20891"
	} else {
		if !strings.HasPrefix(listen, ":") {
			listen = ":" + listen
		}
	}
	log.Printf("echo_ip_httpsvr is listening on port %v", listen)
	log.Printf("examples:")
	log.Printf("http://127.0.0.1%v/", listen)
	log.Printf("http://127.0.0.1%v/ip", listen)
	log.Printf("http://127.0.0.1%v/ip/geo/216.58.221.238", listen)
	log.Printf(`curl -X POST --data '{"username": "xyz", "password": "xyz"}' http://127.0.0.1%v/echo`, listen)
	err := server.ListenAndServe(listen)
	if err != nil {
		log.Fatal(err)
	}
	select {}
}
