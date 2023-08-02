package metrics

import (
	"sync"
)

type MetricsRepository interface {
	Save(metrics Metrics) error
	GetAll() []Metrics
}

type MemoryMetricsRepository struct {
	mutex   sync.Mutex
	metrics []Metrics
}

func NewMemoryMetricsRepository() *MemoryMetricsRepository {
	return &MemoryMetricsRepository{
		mutex:   sync.Mutex{},
		metrics: []Metrics{},
	}
}

func (r *MemoryMetricsRepository) Save(metrics Metrics) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.metrics = append(r.metrics, metrics)
	return nil
}

func (r *MemoryMetricsRepository) GetAll() []Metrics {
	r.mutex.Lock()
	r.mutex.Unlock()
	metricsCopy := make([]Metrics, len(r.metrics))
	copy(metricsCopy, r.metrics)
	return metricsCopy
}
