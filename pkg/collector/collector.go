package collector

import (
	"sync"
	"time"

	"github.com/efremov-it/backup_exporter/pkg/config"
	"github.com/prometheus/client_golang/prometheus"
)

type CollectorMetrics struct {
	mutex                sync.RWMutex
	resultCreate         Result
	resultRetain         Result
	config               *config.Config
	BackupCreateStatus   *prometheus.Desc
	BackupCreateDuration *prometheus.Desc
	BackupRetainStatus   *prometheus.Desc
	BackupRetainDuration *prometheus.Desc
}

type Result struct {
	BackupStartTime time.Time
	BackupDuration  float64
	BackupStatus    int
}

func NewCollector(config *config.Config) *CollectorMetrics {
	labelNames := []string{"backup_type", "project_name", "instance_name"}
	return &CollectorMetrics{
		BackupCreateStatus: prometheus.V2.NewDesc("backup_create_status",
			"Show backup status",
			prometheus.UnconstrainedLabels(labelNames),
			prometheus.Labels(nil)),
		BackupCreateDuration: prometheus.V2.NewDesc("backup_create_duration",
			"Show backup duration",
			prometheus.UnconstrainedLabels(labelNames),
			prometheus.Labels(nil)),
		BackupRetainStatus: prometheus.V2.NewDesc("backup_retain_status",
			"Show backup retain status (only for postgresql)",
			prometheus.UnconstrainedLabels(labelNames),
			prometheus.Labels(nil)),
		BackupRetainDuration: prometheus.V2.NewDesc("backup_retain_duration",
			"Show backup retain duration (only for postgresql)",
			prometheus.UnconstrainedLabels(labelNames),
			prometheus.Labels(nil)),
		config: config,
	}
}

func (collector *CollectorMetrics) Describe(ch chan<- *prometheus.Desc) {
	ds := []*prometheus.Desc{
		collector.BackupCreateStatus,
		collector.BackupCreateDuration,
		collector.BackupRetainStatus,
		collector.BackupRetainDuration,
	}

	for _, d := range ds {
		ch <- d
	}
}

func (c *CollectorMetrics) Collect(ch chan<- prometheus.Metric) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	ch <- prometheus.NewMetricWithTimestamp(c.resultCreate.BackupStartTime, prometheus.MustNewConstMetric(
		c.BackupCreateDuration,
		prometheus.GaugeValue,
		c.resultCreate.BackupDuration,
		c.config.BackupType, c.config.ProjectName, c.config.Host,
	))
	ch <- prometheus.NewMetricWithTimestamp(c.resultCreate.BackupStartTime, prometheus.MustNewConstMetric(
		c.BackupCreateStatus,
		prometheus.GaugeValue,
		float64(c.resultCreate.BackupStatus),
		c.config.BackupType, c.config.ProjectName, c.config.Host,
	))
	ch <- prometheus.NewMetricWithTimestamp(c.resultRetain.BackupStartTime, prometheus.MustNewConstMetric(
		c.BackupRetainDuration,
		prometheus.GaugeValue,
		c.resultRetain.BackupDuration,
		c.config.BackupType, c.config.ProjectName, c.config.Host,
	))
	ch <- prometheus.NewMetricWithTimestamp(c.resultRetain.BackupStartTime, prometheus.MustNewConstMetric(
		c.BackupRetainStatus,
		prometheus.GaugeValue,
		float64(c.resultRetain.BackupStatus),
		c.config.BackupType, c.config.ProjectName, c.config.Host,
	))
}

func (c *CollectorMetrics) SetResultCreate(result Result) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.resultCreate = result
}

func (c *CollectorMetrics) SetResultRetain(result Result) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.resultRetain = result
}
