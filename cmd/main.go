package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/luka2220/tools/rate-limiter/pkg/tokenbucket"
)

var (
	logger = log.New(os.Stdout, "[SERVER]: ", log.LstdFlags)
)

func unlimitedRoute(w http.ResponseWriter, req *http.Request) {
	ip := req.RemoteAddr
	logger.Printf("unlimted route requested by %s\n", ip)
	io.WriteString(w, "unlimited request route...")
}

func limitedRoute(w http.ResponseWriter, req *http.Request) {
	ip := req.RemoteAddr

	// Get the token bucket for the current IP
	bucket := tokenbucket.GetBucket(ip)
	type response struct {
		Message    string `json:"message"`
		Ip         string `json:"ip"`
		BucketSize int    `json:"tokenBucketSize"`
	}
	respUnserialized := &response{
		Message:    "Limited route requested from server...",
		Ip:         bucket.IpAdder,
		BucketSize: bucket.GetBucketSize(),
	}
	respSerialized, err := json.Marshal(respUnserialized)
	if err != nil {
		panic(fmt.Sprintf("Error serializing response struct: %v\n", err))
	}

	// logger.Printf("limited route requested by %s\n", ip)
	// io.WriteString(w, "limited route request ...\n")
	_, err = w.Write(respSerialized)
	if err != nil {
		panic(fmt.Sprintf("Error writting serialized struct to response writer: %v\n", err))
	}
}

func main() {
	unlim := unlimitedRoute
	lim := limitedRoute

	http.HandleFunc("/unlimited", unlim)
	http.HandleFunc("/limited", lim)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
