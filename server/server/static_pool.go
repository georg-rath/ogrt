package server

import (
	"sync"
	"sync/atomic"

	"github.com/rcrowley/go-metrics"
)

type StaticPool struct {
	pool     sync.Pool
	size     int
	inFlight int64
}

func (s StaticPool) Get() []byte {
	atomic.AddInt64(&s.inFlight, 1)
	return s.pool.Get().([]byte)
}

func (s StaticPool) Put(buf []byte) {
	if cap(buf) != s.size {
		panic("buffer with wrong size returned to pool")
	}
	atomic.AddInt64(&s.inFlight, -1)
	s.pool.Put(buf)
}

func (s StaticPool) InFlight() int64 {
	inFlight := atomic.LoadInt64(&s.inFlight)
	return inFlight
}

func NewStaticPool(name string, size int) (s *StaticPool) {
	s = &StaticPool{
		pool: sync.Pool{
			New: func() interface{} { return make([]byte, size) },
		},
		size: size,
	}
	metrics.NewRegisteredFunctionalGauge(name+"_inflight", metrics.DefaultRegistry, func() int64 {
		return s.InFlight()
	})

	return
}
