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

func (ips *IpAddressStore) getBucket(ip string) tokenbucket.Bucket {
	bucket, ok := ips.ipAddress[ip]
	if ok {
		// TODO: since the ip already exists, call token remove from tokenBucket
		return bucket
	} else {
		// TODO: since a new ip address is added store the ip and it's new bucket to the map store
		ips.ipAddress[ip] = *tokenbucket.NewBucket()
		return ips.ipAddress[ip]
	}
}

func (ips *IpAddressStore) unlimitedRoute(w http.ResponseWriter, req *http.Request) {
	ip := req.RemoteAddr
	logger.Printf("unlimted route requested by %s\n", ip)
	io.WriteString(w, "unlimited request route...")
}

func (ips *IpAddressStore) limitedRoute(w http.ResponseWriter, req *http.Request) {
	ip := req.RemoteAddr
	bucket := ips.getBucket(ip)

	logger.Printf("limited route requested by %s\n", ip)
	logger.Printf("current ip bucket has %d tokens", bucket.GetBucketSize())
	io.WriteString(w, "limited request route...")
}

func main() {
	ipStorage := newIpAddressStore()

	unlim := ipStorage.unlimitedRoute
	lim := ipStorage.limitedRoute

	http.HandleFunc("/unlimited", unlim)
	http.HandleFunc("/limited", lim)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
