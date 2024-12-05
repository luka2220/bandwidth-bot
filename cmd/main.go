package main

import (
	"encoding/json"
	"fmt"
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

	type response struct {
		Message string `json:"message"`
		Ip      string `json:"ip"`
	}

	respSerialized, err := json.Marshal(&response{
		Message: "Unlimited route requested from the server...",
		Ip:      ip,
	})
	if err != nil {
		panic(fmt.Sprintf("Error serializing struct to json response: %v\n", err))
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(respSerialized)
	if err != nil {
		panic(fmt.Sprintf("Error writting serialized struct to response writer: %v\n", err))
	}
}

func limitedRoute(w http.ResponseWriter, req *http.Request) {
	ip := req.RemoteAddr

	bucket := tokenbucket.GetIpAdderBucket(ip)
	type response struct {
		Message string `json:"message"`
		Ip      string `json:"ip"`
	}

	var respUnserialized *response

	switch bucket.GetHTTPStatus() {
	case 429:
		respUnserialized = &response{
			Message: "The client has sent too many requests in a given amount of time",
			Ip:      ip,
		}

	case 200:
		respUnserialized = &response{
			Message: "Limited route requested from server...",
			Ip:      ip,
		}
	}

	respSerialized, err := json.Marshal(respUnserialized)
	if err != nil {
		panic(fmt.Sprintf("Error serializing response struct: %v\n", err))
	}

	logger.Printf("limited route requested by %s\n", ip)

	w.WriteHeader(bucket.GetHTTPStatus())
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
