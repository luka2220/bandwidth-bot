package bandwidthbot

import (
	"log"
	"os"
	"sync"
	"time"
)

var (
	fixedWindowStore      = make(map[string]*FixedWindow)
	loggerFW              = log.New(os.Stdout, "[SERVER-FWC]: ", log.LstdFlags)
	fixedWindowStoreMutex sync.Mutex
)

type FixedWindow struct {
	ipAddress   string    // Ip address of the current window
	windowSize  float64   // Time limit of window in seconds
	requests    int       // Counter for requests
	maxRequests int       // Maximum allowed requests within the window
	startTime   time.Time // Start time of the current window
}

func (fw *FixedWindow) RunFixedWindow(ipAdder string) (int, string) {
	fixedWindowStoreMutex.Lock()
	defer fixedWindowStoreMutex.Unlock()

}
