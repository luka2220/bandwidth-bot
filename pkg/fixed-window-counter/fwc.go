package fixedwindowcounter

import "time"

type FixedWindow struct {
	threshold float64   // Limit where if exceded the requests will be discarded
	counter   float64   // Amount to increase window time per incoming request
	window    time.Time // Current time of window
}

func UseFixedWindow() *FixedWindow {
	return &FixedWindow{}
}
