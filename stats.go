package gcache

import (
	"sync/atomic"
)

type statsAccessor interface {
	HitCount() uint64
	MissCount() uint64
	ExpireCount() uint64
	LookupCount() uint64
	HitRate() float64
}

// statistics
type stats struct {
	hitCount      uint64
	missCount     uint64
	expireCount   uint64
	hitMetrics    MetricsFunc
	missMetrics   MetricsFunc
	expireMetrics MetricsFunc
}

// increment hit count
func (st *stats) IncrHitCount() uint64 {
	atomic.AddUint64(&st.hitCount, 1)
	if st.hitMetrics != nil {
		st.hitMetrics(st.HitCount(), 1)
	}
	return st.HitCount()
}

// increment miss count
func (st *stats) IncrMissCount() uint64 {
	atomic.AddUint64(&st.missCount, 1)
	if st.missMetrics != nil {
		st.missMetrics(st.MissCount(), 1)
	}
	return st.MissCount()
}

// increment expire count
func (st *stats) IncrExpireCount() uint64 {
	atomic.AddUint64(&st.expireCount, 1)
	if st.expireMetrics != nil {
		st.expireMetrics(st.ExpireCount(), 1)
	}
	return st.ExpireCount()
}

// HitCount returns hit count
func (st *stats) HitCount() uint64 {
	return atomic.LoadUint64(&st.hitCount)
}

// MissCount returns miss count
func (st *stats) MissCount() uint64 {
	return atomic.LoadUint64(&st.missCount)
}

// ExpireCount returns expire count
func (st *stats) ExpireCount() uint64 {
	return atomic.LoadUint64(&st.expireCount)
}

// LookupCount returns lookup count
func (st *stats) LookupCount() uint64 {
	return st.HitCount() + st.MissCount() + st.ExpireCount()
}

// HitRate returns rate for cache hitting
func (st *stats) HitRate() float64 {
	hc, mc, ec := st.HitCount(), st.MissCount(), st.ExpireCount()
	total := hc + mc + ec
	if total == 0 {
		return 0.0
	}
	return float64(hc+ec) / float64(total)
}
