package bandwidthbot_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/luka2220/bandwidthbot"
)

func limited(w http.ResponseWriter, req *http.Request) {
	ip := req.RemoteAddr

	fwc := bandwidthbot.RunFixedWindow(ip)
	serverResponseCode := fwc.GetHTTPStatus()

	response := struct {
		Message string `json:"message"`
		Status  int    `json:"status"`
		Ip      string `json:"ip"`
	}{
		Message: "rate limiting api route",
		Status:  serverResponseCode,
		Ip:      ip,
	}

	byteResponse, err := json.Marshal(response)
	if err != nil {
		panic(fmt.Sprintf("Error serializing data: %s", err))
	}
	w.WriteHeader(serverResponseCode)
	_, err = w.Write(byteResponse)
	if err != nil {
		panic(fmt.Sprintf("Error writting data to http client: %s", err))
	}
}

func startRateLimiterServer() {
	lim := limited

	http.HandleFunc("/limited", lim)
	http.ListenAndServe(":8080", nil)
}

func TestBandwidthBot(t *testing.T) {
	startRateLimiterServer()
}
