package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/luka2220/bandwidthbot"
)

var (
	logger = log.New(os.Stdout, "[SERVER]: ", log.LstdFlags)
)

const (
	TOKEN_BUCKET         = "token-bucket"
	FIXED_WINDOW_COUNTER = "fixed-window-counter"
)

type response struct {
	Message string `json:"message"`
	Ip      string `json:"ip"`
}

type rateLimiter struct {
	name string
}

func newRateLimiter(name string) *rateLimiter {
	return &rateLimiter{
		name,
	}
}

func (r rateLimiter) unlimited(w http.ResponseWriter, req *http.Request) {
	ip := req.RemoteAddr
	logger.Printf("unlimted route requested by %s\n", ip)

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

func (r rateLimiter) limited(w http.ResponseWriter, req *http.Request) {
	ip := req.RemoteAddr

	var serverResponseCode = 200

	switch r.name {
	case TOKEN_BUCKET:
		bucket := bandwidthbot.InitializeTokenBucket(ip)
		serverResponseCode = bucket.GetHTTPStatus()
	case FIXED_WINDOW_COUNTER:
		fwc := bandwidthbot.InitializeFixedWindow(ip)
		serverResponseCode = fwc.GetHTTPStatus()
	}

	var respUnserialized *response

	switch serverResponseCode {
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

	w.WriteHeader(serverResponseCode)
	_, err = w.Write(respSerialized)
	if err != nil {
		panic(fmt.Sprintf("Error writting serialized struct to response writer: %v\n", err))
	}
}

func main() {
	rl := newRateLimiter(FIXED_WINDOW_COUNTER)

	unlim := rl.unlimited
	lim := rl.limited

	http.HandleFunc("/unlimited", unlim)
	http.HandleFunc("/limited", lim)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
