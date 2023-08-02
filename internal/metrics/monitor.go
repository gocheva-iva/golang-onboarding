package metrics

import (
	"fmt"

	"github.com/gocheva-iva/task1/internal/logging"
	"github.com/labstack/echo/v4"
	"github.com/robfig/cron/v3"
)

type MetricsMonitor struct {
	Repository MetricsRepository
	Collector  MetricsCollector
	Cron       *cron.Cron
	threshold  int
}

func NewMetricsMonitor(repository MetricsRepository, collector MetricsCollector, threshold int) *MetricsMonitor {
	return &MetricsMonitor{
		Repository: repository,
		Collector:  collector,
		threshold:  threshold,
	}
}

func (m *MetricsMonitor) Invoke(ctx echo.Context) error {
	metrics, err := m.Collector.Collect(ctx)
	if err != nil {
		return fmt.Errorf("error collecting metrics %w", err)
	}
	m.LogMetrics(metrics)
	err = m.Repository.Save(metrics)
	if err != nil {
		return fmt.Errorf("error saving metrics %w", err)
	}
	return nil
}

func (m *MetricsMonitor) LogMetrics(metrics Metrics) {
	if metrics.AboveThreshold(m.threshold) {
		json, err := metrics.ToJSONString()
		if err != nil {
			logging.GetLogger().Error("error marshalling metrics: %w", err)
			return
		}
		logging.GetLogger().FileLog("repository.txt", json)
	}
}
