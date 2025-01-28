package metrics

import (
	"math/rand/v2"
	"runtime"
	"sync"
	"time"
)

const (
	pollInterval = 2 * time.Second
)

func NewDB() *DB {
	return &DB{
		Metrics: Metrics{},
		mu:      sync.RWMutex{},
	}
}

func (db *DB) RunUpdates() {
	go func() {
		for {
			db.poll()
			time.Sleep(pollInterval)
		}
	}()
}

func (db *DB) poll() {
	ms := runtime.MemStats{}
	runtime.ReadMemStats(&ms)
	randomValue := Gauge(rand.Float64())

	db.mu.Lock()
	defer db.mu.Unlock()

	db.Metrics.Alloc = Gauge(ms.Alloc)
	db.Metrics.BuckHashSys = Gauge(ms.BuckHashSys)
	db.Metrics.Frees = Gauge(ms.Frees)
	db.Metrics.GCCPUFraction = Gauge(ms.GCCPUFraction)
	db.Metrics.GCSys = Gauge(ms.GCSys)
	db.Metrics.HeapAlloc = Gauge(ms.HeapAlloc)
	db.Metrics.HeapIdle = Gauge(ms.HeapIdle)
	db.Metrics.HeapInuse = Gauge(ms.HeapInuse)
	db.Metrics.HeapObjects = Gauge(ms.HeapObjects)
	db.Metrics.HeapReleased = Gauge(ms.HeapReleased)
	db.Metrics.HeapSys = Gauge(ms.HeapSys)
	db.Metrics.LastGC = Gauge(ms.LastGC)
	db.Metrics.Lookups = Gauge(ms.Lookups)
	db.Metrics.MCacheInuse = Gauge(ms.MCacheInuse)
	db.Metrics.MCacheSys = Gauge(ms.MCacheSys)
	db.Metrics.MSpanInuse = Gauge(ms.MSpanInuse)
	db.Metrics.MSpanSys = Gauge(ms.MSpanSys)
	db.Metrics.Mallocs = Gauge(ms.Mallocs)
	db.Metrics.NextGC = Gauge(ms.NextGC)
	db.Metrics.NumForcedGC = Gauge(ms.NumForcedGC)
	db.Metrics.NumGC = Gauge(ms.NumGC)
	db.Metrics.OtherSys = Gauge(ms.OtherSys)
	db.Metrics.PauseTotalNs = Gauge(ms.PauseTotalNs)
	db.Metrics.StackInuse = Gauge(ms.StackInuse)
	db.Metrics.StackSys = Gauge(ms.StackSys)
	db.Metrics.Sys = Gauge(ms.Sys)
	db.Metrics.TotalAlloc = Gauge(ms.TotalAlloc)

	db.Metrics.RandomValue = randomValue
	db.Metrics.PollCounter++
}

func (db *DB) GetMetrics() Metrics {
	db.mu.RLock()
	defer db.mu.RUnlock()
	return db.Metrics
}
