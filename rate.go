package rhombus

import (
	"time"
)

// Convenience type for sleeping in a loop at a
// specified rate
type Rate struct{
	dur time.Duration
	lastTime time.Time
}

// create a new Rate object
// freq: desired rate to run at in Hz
func NewRate(freq uint) *Rate {
	return &Rate{
		dur: time.Duration(1e9 / freq) * time.Nanosecond,
		lastTime: time.Now(),
	}
}

// attempt to sleep at a specified rate
func (r *Rate) Sleep() {
	currTime := time.Now()
	elapsed := currTime.Sub(r.lastTime)
	time.Sleep(r.dur - elapsed)

	// detect time jumping backwards of forwards
	if r.lastTime.After(currTime) ||
		currTime.Sub(r.lastTime) > r.dur * 2 {
		r.lastTime = currTime
	} else {
		r.lastTime = r.lastTime.Add(r.dur)
	}
}
