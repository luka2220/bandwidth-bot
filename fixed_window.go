package bandwidthbot

import (
	"log"
	"net/http"
	"os"
	"time"
)

var (
	fixedWindowStore = make(map[string]*FixedWindow)
	loggerFW         = log.New(os.Stdout, "[SERVER-FWC]: ", log.LstdFlags)
)

type FixedWindow struct {
	ipAddress  string    // Ip address of the current fixed window
	windowSize float64   // Limit where if exceded the requests will be discarded
	counter    float64   // Amount to increase window time per incoming request
	startTime  time.Time // Start time of the window
	httpStatus int       // Current http status
}

func InitializeFixedWindow(ipAdder string) *FixedWindow {
	window, ok := fixedWindowStore[ipAdder]
	if ok {
		elpased_time := time.Now().Sub(window.startTime)
		loggerFW.Printf("Elpased time = %.2f\n", elpased_time.Seconds())

		if elpased_time.Seconds() > window.windowSize {
			fixedWindowStore[ipAdder].httpStatus = http.StatusTooManyRequests
		}

		return window
	}

	newFixedWindow := &FixedWindow{
		ipAddress:  ipAdder,
		windowSize: 60,
		counter:    5.00,
		startTime:  time.Now(),
		httpStatus: 200,
	}

	fixedWindowStore[ipAdder] = newFixedWindow
	loggerFW.Printf("%s created in fixedWindowCounter store\n", newFixedWindow.ipAddress)

	return newFixedWindow
}

func (fw *FixedWindow) increaseWindowCounter() {

}

func (fw *FixedWindow) GetHTTPStatus() int {
	return fw.httpStatus

}
