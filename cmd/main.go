package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/luka2220/tools/rate-limiter/pkg/tokenbucket"
)

var (
	logger = log.New(os.Stdout, "[SERVER]: ", log.LstdFlags)
)


type IpAddressStore struct {
	ipAddress map[string]tokenbucket.Bucket
}

func newIpAddressStore() *IpAddressStore {
	return &IpAddressStore{
		make(map[string]tokenbucket.Bucket),
	}
}

func unlimitedRoute(w http.ResponseWriter, req *http.Request) {
	clientHost := req.RemoteAddr
	logger.Printf("unlimted route requested by %s\n", clientHost)
	io.WriteString(w, "unlimited request route...")
}

func limitedRoute(w http.ResponseWriter, req *http.Request) {
	ip := req.RemoteAddr
	logger.Printf("limited route requested by %s\n", ip)
	io.WriteString(w, "limited request route...")
}

func main() {
	unlim := unlimitedRoute
	lim := limitedRoute

	http.HandleFunc("/unlimited", unlim)
	http.HandleFunc("/limited", lim)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
