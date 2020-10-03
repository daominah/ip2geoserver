package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"

	"github.com/daominah/echo_ip_httpsvr/ip2geo"
	"github.com/mywrap/gofast"
)

func main() {
	handler := http.NewServeMux()
	handler.HandleFunc("/", func() http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			log.Printf("_____________________________________________\n")
			log.Printf("begin HandleFunc echo %v\n", r.RemoteAddr)
			rDump, err := httputil.DumpRequest(r, true)
			if err != nil {
				w.WriteHeader(500)
				w.Write([]byte(fmt.Sprintf("error DumpRequest: %v", err)))
				return
			}
			w.Write([]byte(fmt.Sprintf("echo req from %v:\n\n", r.RemoteAddr)))
			w.Write(rDump)
			log.Printf("respond to %v: %s\n", r.RemoteAddr, rDump)
			log.Printf("end HandleFunc echo %v\n", r.RemoteAddr)
			log.Printf("_____________________________________________\n")
			return
		}
	}())
	handler.HandleFunc("/ip", func() http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ip := ip2geo.GetIpFromAddress(r.RemoteAddr)
			w.Write([]byte(fmt.Sprintf("%v", ip)))
			return
		}
	}())
	handler.HandleFunc("/ip/geo", func() http.HandlerFunc {
		prjRoot, err := gofast.GetProjectRootPath()
		if err != nil {
			log.Fatal(err)
		}
		geoReader, err := ip2geo.NewReader(prjRoot + "/ip2geo")
		if err != nil {
			log.Fatal(err)
		}
		return func(w http.ResponseWriter, r *http.Request) {
			ip := ip2geo.GetIpFromAddress(r.RemoteAddr)
			geo, err := geoReader.ReadIPInfo(ip)
			if err != nil {
				w.WriteHeader(400)
				w.Write([]byte(err.Error()))
			}
			beauty, err := json.MarshalIndent(geo, "", "\t")
			if err != nil {
				w.WriteHeader(500)
				w.Write([]byte(err.Error()))
			}
			w.Write(beauty)
			return
		}
	}())

	listen := os.Getenv("LISTENING_PORT")
	if listen == "" {
		listen = ":20891"
	} else {
		if !strings.HasPrefix(listen, ":") {
			listen = ":" + listen
		}
	}
	server := &http.Server{Addr: listen, Handler: handler}
	log.Println("echo server listening on port ", server.Addr)
	log.Printf("example: http://127.0.0.1%v/", server.Addr)
	log.Printf("example: http://127.0.0.1%v/ip", server.Addr)
	log.Printf("example: http://127.0.0.1%v/ip/geo", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

	// curl 127.0.0.1:20891 --request POST --data '{"name":"Tung"}'
	// 127.0.0.1:20891?a=5

	select {}
}
