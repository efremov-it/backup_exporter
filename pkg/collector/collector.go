package collector

import (
	"sync"
	"time"

	"github.com/efremov-it/backup_exporter/pkg/config"
	"github.com/prometheus/client_golang/prometheus"
)

type CollectorMetrics struct {
	mutex          sync.RWMutex
	result         Result
	config         *config.Config
	BackupStatus   *prometheus.Desc
	BackupDuration *prometheus.Desc
	BackupNextTime *prometheus.Desc
}

type Result struct {
	BackupStartTime time.Time
	BackupDuration  float64
	BackupStatus    int
}

func NewCollector(config *config.Config) *CollectorMetrics {
	labelNames := []string{"backup_type", "project_name", "instance_name"}
	return &CollectorMetrics{
		BackupStatus: prometheus.V2.NewDesc("backup_status",
			"Show backup status",
			prometheus.UnconstrainedLabels(labelNames),
			prometheus.Labels(nil)),
		BackupDuration: prometheus.V2.NewDesc("backup_duration",
			"Show backup duration",
			prometheus.UnconstrainedLabels(labelNames),
			prometheus.Labels(nil)),
		BackupNextTime: prometheus.V2.NewDesc("backup_next_time",
			"Show when the next backup will be created",
			prometheus.UnconstrainedLabels(labelNames),
			prometheus.Labels(nil)),
		config: config,
	}
}

func (collector *CollectorMetrics) Describe(ch chan<- *prometheus.Desc) {
	ds := []*prometheus.Desc{
		collector.BackupStatus,
		collector.BackupDuration,
		collector.BackupNextTime,
	}

	for _, d := range ds {
		ch <- d
	}
}

func (c *CollectorMetrics) Collect(ch chan<- prometheus.Metric) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	ch <- prometheus.NewMetricWithTimestamp(c.result.BackupStartTime, prometheus.MustNewConstMetric(
		c.BackupDuration,
		prometheus.GaugeValue,
		c.result.BackupDuration,
		c.config.BackupType, c.config.ProjectName, c.config.Host,
	))
	ch <- prometheus.NewMetricWithTimestamp(c.result.BackupStartTime, prometheus.MustNewConstMetric(
		c.BackupStatus,
		prometheus.GaugeValue,
		float64(c.result.BackupStatus),
		c.config.BackupType, c.config.ProjectName, c.config.Host,
	))
}

func (c *CollectorMetrics) SetResult(result Result) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.result = result
}
