package gcache

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestStats(t *testing.T) {
	var cases = []struct {
		hit  int
		miss int
		rate float64
	}{
		{3, 1, 0.75},
		{0, 1, 0.0},
		{3, 0, 1.0},
		{0, 0, 0.0},
	}

	for _, cs := range cases {
		st := &stats{}
		for i := 0; i < cs.hit; i++ {
			st.IncrHitCount()
		}
		for i := 0; i < cs.miss; i++ {
			st.IncrMissCount()
		}
		if rate := st.HitRate(); rate != cs.rate {
			t.Errorf("%v != %v", rate, cs.rate)
		}
	}
}

func getter(key interface{}) (interface{}, error) {
	return key, nil
}

func TestCacheStats(t *testing.T) {
	var cases = []struct {
		builder func() Cache
		rate    float64
	}{
		{
			builder: func() Cache {
				cc := New(32).Simple().Build()
				cc.Set(0, 0)
				cc.Get(0)
				cc.Get(1)
				return cc
			},
			rate: 0.5,
		},
		{
			builder: func() Cache {
				cc := New(32).LRU().Build()
				cc.Set(0, 0)
				cc.Get(0)
				cc.Get(1)
				return cc
			},
			rate: 0.5,
		},
		{
			builder: func() Cache {
				cc := New(32).LFU().Build()
				cc.Set(0, 0)
				cc.Get(0)
				cc.Get(1)
				return cc
			},
			rate: 0.5,
		},
		{
			builder: func() Cache {
				cc := New(32).ARC().Build()
				cc.Set(0, 0)
				cc.Get(0)
				cc.Get(1)
				return cc
			},
			rate: 0.5,
		},
		{
			builder: func() Cache {
				cc := New(32).
					Simple().
					LoaderFunc(getter).
					Build()
				cc.Set(0, 0)
				cc.Get(0)
				cc.Get(1)
				return cc
			},
			rate: 0.5,
		},
		{
			builder: func() Cache {
				cc := New(32).
					LRU().
					LoaderFunc(getter).
					Build()
				cc.Set(0, 0)
				cc.Get(0)
				cc.Get(1)
				return cc
			},
			rate: 0.5,
		},
		{
			builder: func() Cache {
				cc := New(32).
					LFU().
					LoaderFunc(getter).
					Build()
				cc.Set(0, 0)
				cc.Get(0)
				cc.Get(1)
				return cc
			},
			rate: 0.5,
		},
		{
			builder: func() Cache {
				cc := New(32).
					ARC().
					LoaderFunc(getter).
					Build()
				cc.Set(0, 0)
				cc.Get(0)
				cc.Get(1)
				return cc
			},
			rate: 0.5,
		},
	}

	for i, cs := range cases {
		cc := cs.builder()
		if rate := cc.HitRate(); rate != cs.rate {
			t.Errorf("case-%v: %v != %v", i, rate, cs.rate)
		}
	}
}

func TestStatsMetrics(t *testing.T) {
	var testCounter int64
	var hitCounter int64
	var missCounter int64
	var expireCounter int64

	expire := 200 * time.Millisecond
	counter := 1000

	cache := New(1).Simple().
		LoaderExpireFunc(func(key interface{}) (interface{}, *time.Duration, error) {
			time.Sleep(10 * time.Millisecond)
			return atomic.AddInt64(&testCounter, 1), &expire, nil
		}, true).
		HitMetricsFunc(func(u uint64, i int64) {
			atomic.AddInt64(&hitCounter, i)
		}).MissMetricsFunc(func(u uint64, i int64) {
			atomic.AddInt64(&missCounter, i)
		}).ExpireMetricsFunc(func(u uint64, i int64) {
			atomic.AddInt64(&expireCounter, i)
		}).Build()


	wg := sync.WaitGroup{}
	for i := 0; i < counter; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := cache.Get(0)
			if err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
	if missCounter != 1000 {
		t.Errorf("missCounter = %v, should be %v", missCounter, cache.MissCount())
	}

	for i := 0; i < counter; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := cache.Get(0)
			if err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
	if hitCounter != 1000 {
		t.Errorf("hitCounter = %v, should be %v", hitCounter, cache.HitCount())
	}

	time.Sleep(expire)
	for i := 0; i < counter; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := cache.Get(0)
			if err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
	if expireCounter != 1000 {
		t.Errorf("expireCounter = %v, should be %v", expireCounter, cache.ExpireCount())
	}
}
