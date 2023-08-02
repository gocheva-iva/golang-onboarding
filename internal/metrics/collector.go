package metrics

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gocheva-iva/task1/internal/logging"
	"github.com/labstack/echo/v4"
)

type MetricsCollector interface {
	Collect(echo.Context) (Metrics, error)
}

type MetricsDummyGateway struct {
	Client *http.Client
}

const (
	endpoint = "http://localhost:1323/"
	Timeout  = 4
)

func NewMetricsDummyGateway(client *http.Client) *MetricsDummyGateway {
	return &MetricsDummyGateway{
		Client: client,
	}
}

func (m *MetricsDummyGateway) Collect(rootCtx echo.Context) (Metrics, error) {
	ctx, cancel := context.WithTimeout(rootCtx.Request().Context(), time.Second*time.Duration(Timeout))
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return Metrics{}, fmt.Errorf("error creating request: %w", err)
	}
	resp, err := m.Client.Do(req)
	if err != nil {
		logging.GetLogger().Error("error sending request", err)
		return Metrics{}, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	var metrics Metrics
	if err := json.NewDecoder(resp.Body).Decode(&metrics); err != nil {
		return Metrics{}, fmt.Errorf("error decoding response body: %w", err)
	}
	return metrics, nil
}
