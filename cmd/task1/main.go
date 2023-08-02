package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gocheva-iva/task1/internal/logging"
	"github.com/gocheva-iva/task1/internal/metrics"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func collectMetrics(repo *metrics.MemoryMetricsRepository) echo.HandlerFunc {
	return func(c echo.Context) error {
		metrics := repo.GetAll()
		err := c.JSON(http.StatusOK, metrics)
		if err != nil {
			return fmt.Errorf("error marshalling metrics: %w", err)
		}
		return nil
	}
}

func displayMetrics(monitor *metrics.MetricsMonitor) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := monitor.Invoke(c)
		if err != nil {
			logging.GetLogger().Error("error invoking monitor: %w", err)
		}
		return nil
	}
}

func main() {
	notMagicNum := 10
	repo := metrics.NewMemoryMetricsRepository()
	collector := metrics.NewMetricsDummyGateway(http.DefaultClient)
	monitor := metrics.NewMetricsMonitor(repo, collector, notMagicNum)

	e := echo.New()
	defer e.Close()
	logging.GetLogger().Info("Starting server...")
	e.GET("/data", collectMetrics(repo))
	e.GET("/metrics", displayMetrics(monitor))
	err := godotenv.Load()
	if err != nil {
		logging.GetLogger().Error("Error loading .env file", err)
		return
	}
	port := os.Getenv("PORT")
	if port == "" {
		logging.GetLogger().Error("$PORT not set")
		return
	}

	err = e.Start(":" + port)
	if err != nil {
		logging.GetLogger().Error("Error starting Echo server:", err)
		return
	}
}
