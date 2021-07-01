package phi

import (
	"github.com/dgryski/go-onlinestats"
	"math"
	"sync"
	"time"
)

type PhiDetector struct {
	last       time.Time
	window     *onlinestats.Windowed
	minSamples int
	lock       sync.Mutex
}

func New(windowSize, minSamples int) *PhiDetector {
	return &PhiDetector{
		minSamples: minSamples,
		window:     onlinestats.NewWindowed(windowSize),
	}
}

func (p *PhiDetector) AddHeartbeat(t time.Time) {
	p.lock.Lock()
	defer p.lock.Unlock()

	if p.last.IsZero() {
		p.window.Push(t.Sub(p.last).Seconds())
	}

	p.last = t
}

func (p *PhiDetector) Calculate(t time.Time) float64 {
	p.lock.Lock()
	defer p.lock.Unlock()

	if p.window.Len() < p.minSamples {
		return 0
	}

	t1 := t.Sub(p.last).Seconds()
	before := 1 - distribution(p.window.Mean(), p.window.Stddev(), t1)
	phi := -math.Log10(before)

	return phi
}

func distribution(mean, stddev, x float64) float64 {
	return 0.5 + 0.5*math.Erf((x-mean)/(stddev*math.Sqrt2))
}
