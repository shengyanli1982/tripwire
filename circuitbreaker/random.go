package circuitbreaker

import (
	"math/rand"
	"sync"
	"time"
)

type SafeRandom struct {
	// rand.New(...) returns a non thread safe object
	r    *rand.Rand
	lock sync.Mutex
}

// NewSafeRandom returns a Proba.
func NewSafeRandom() *SafeRandom {
	return &SafeRandom{
		r: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// LessThan checks if true on given probability.
func (p *SafeRandom) LessThan(v float64) bool {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.r.Float64() < v
}

// GreaterThan checks if true on given probability.
func (p *SafeRandom) GreaterThan(v float64) bool {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.r.Float64() >= v
}

func (p *SafeRandom) Float64() float64 {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.r.Float64()
}
