package bandwidthbot

import (
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	fixedWindowStore      = make(map[string]*FixedWindow)
	loggerFW              = log.New(os.Stdout, "[SERVER-FWC]: ", log.LstdFlags)
	fixedWindowStoreMutex sync.Mutex
	expiryDuration   = 30 * time.Second
)

type FixedWindow struct {
	ipAddress   string    // Ip address of the current window
	windowSize  float64   // Time limit of window in seconds
	requests    int       // Counter for requests
	maxRequests int       // Maximum allowed requests within the window
	startTime   time.Time // Start time of the current window
	lastRequest time.Time	// Time since the ip last made a request
}

// RunFixedWindow executes the fixed window counter
// algorithm for each incoming ip and returns the
// current status code i.e 420 for too many requests
func RunFixedWindow(ip string) int {
	fixedWindowStoreMutex.Lock()
	defer fixedWindowStoreMutex.Unlock()

	removeExpiredEntries()

	window, exists := fixedWindowStore[ip]; if ! exists {
		fixedWindowStore[ip] = &FixedWindow{
			ipAddress: ip,
			windowSize: 60.00,
			requests: 1,
			maxRequests: 5,
			startTime: time.Now(),
			lastRequest: time.Now(),
		}

		loggerFW.Printf("Window started for %s\n", ip)
		return http.StatusOK
	}

	window.lastRequest = time.Now()

	if time.Since(window.startTime).Seconds() > window.windowSize {
		window.startTime = time.Now()
		window.requests = 0
		loggerFW.Printf("Fixed window reset for %s\n", ip)
	}

	if window.requests < window.maxRequests {
		window.requests++
		loggerFW.Printf("Request allowed for %s (count: %d)\n", ip, window.requests)
		return http.StatusOK
	}

	loggerFW.Printf("Request denied for %s (too many requests)\n", ip)
	return http.StatusTooManyRequests
}

func removeExpiredEntries() {
	now := time.Now()
	for ip, window := range fixedWindowStore {
		if now.Sub(window.lastRequest) > expiryDuration {
			delete(fixedWindowStore, ip)
			loggerFW.Printf("Expired entry removed for IP: %s\n", ip)
		}
	}
}