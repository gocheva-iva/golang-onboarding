package metrics

import (
	"encoding/json"
	"fmt"

	"github.com/shirou/gopsutil/net"
)

type Metrics struct {
	CPU     float64              `json:"cpu"`
	Memory  uint64               `json:"memory"`
	Network []net.IOCountersStat `json:"network"`
}

func (m Metrics) ToJSONString() (string, error) {
	jsonData, err := json.Marshal(m)
	if err != nil {
		return "", fmt.Errorf("error marshalling metrics: %w", err)
	}
	return string(jsonData) + "\n", nil
}

func (m Metrics) AboveThreshold(threshold int) bool {
	return (m.CPU > float64(threshold)) || (m.Memory > uint64(threshold))
}
