package create

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/efremov-it/backup_exporter/pkg/collector"
	"github.com/efremov-it/backup_exporter/pkg/config"
)

func NewBackupService(config *config.Config, metrics *collector.CollectorMetrics) *BackupService {
	return &BackupService{
		config:  config,
		metrics: metrics,
	}
}

type BackupService struct {
	config  *config.Config
	metrics *collector.CollectorMetrics
}

func (bs *BackupService) Run(ctx context.Context) {
	var command string
	var use_flags []string
	// set backup command and arguments
	if bs.config.BackupType != "clickhouse" {
		// for psql we need to use BackupStorage
		use_flags = append(use_flags, "backup-push")
		if bs.config.BackupType == "psql" {
			use_flags = append(use_flags, bs.config.BackupStorage)
		}
		command = "wal-g"
	} else {
		use_flags = append(use_flags, "create_remote")
		command = "clickhouse-backup"
	}

	if bs.config.ConfigFile != "" {
		use_flags = append(use_flags, "--config", bs.config.ConfigFile)
	}

	// create backup
	n := 0
	for n < 3 {
		result, err := backup(command, use_flags...)
		if err != nil {
			bs.metrics.SetResult(result)
			fmt.Fprint(os.Stdout, []any{"Backup failed\nStart backup one more time\n"}...)
			time.Sleep(5 * time.Second)
		} else {
			n = 3
			print("Backup created\n")
			bs.metrics.SetResult(result)
			return
		}
		n++
	}
}

func (bs *BackupService) Retain(ctx context.Context) {
	var use_flags []string
	// for psql we need to use BackupStorage
	if bs.config.BackupType == "psql" {
		if bs.config.DeleteRetain != "" {
			use_flags = append(use_flags, "delete", "retain", "FULL", bs.config.DeleteRetain, "--confirm")
		}
		use_flags = append(use_flags, bs.config.BackupStorage)
	}

	if bs.config.ConfigFile != "" {
		use_flags = append(use_flags, "--config", bs.config.ConfigFile)
	}

	// Retain backup
	result, err := backup("wal-g", use_flags...)
	if err != nil {
		bs.metrics.SetResult(result)
		fmt.Fprint(os.Stdout, []any{"Retain backup failed\n"}...)
	} else {
		print("Retain backup finished\n")
		bs.metrics.SetResult(result)
		return
	}
}

func backup(command string, args ...string) (result collector.Result, err error) {
	result.BackupStartTime = time.Now()
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()

	defer func() {
		result.BackupDuration = float64(time.Since(result.BackupStartTime).Seconds())
		if err != nil {
			result.BackupStatus = 1
			log.Printf("Error:\n", err, string(output))
		} else {
			result.BackupStatus = 0
			log.Printf("Command Output: %s\n", string(output))
		}
	}()
	return
}
