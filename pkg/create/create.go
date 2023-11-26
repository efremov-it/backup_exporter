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

func (bs *BackupService) Create(ctx context.Context) {
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

	// create backup. if backup failed try do it again 3 times
	for i := 0; i < 3; i++ {
		t, d, s, err := backup(command, use_flags...)
		bs.metrics.SetResultCreate(collector.Result{t, d, s})
		if err != nil {
			fmt.Fprint(os.Stdout, []any{"Backup failed\nStart backup one more time\n"}...)
			time.Sleep(10 * time.Second)
		} else {
			print("Backup created\n")
			return
		}
	}

}

func (bs *BackupService) Retain(ctx context.Context) {
	var use_flags []string
	if bs.config.BackupType == "psql" {
		if bs.config.DeleteRetain != "" {
			use_flags = append(use_flags, "delete", "retain", "FULL", bs.config.DeleteRetain, "--confirm")
		}
		use_flags = append(use_flags, bs.config.BackupStorage)
	} else {
		print("Retain backup support only for postgresql\n")
	}

	if bs.config.ConfigFile != "" {
		use_flags = append(use_flags, "--config", bs.config.ConfigFile)
	}

	// Retain backup
	t, d, s, err := backup("wal-g", use_flags...)
	bs.metrics.SetResultRetain(collector.Result{t, d, s})
	if err != nil {
		fmt.Fprint(os.Stdout, []any{"Retain backup failed\n"}...)
	} else {
		print("Retain backup finished\n")
		return
	}
}

func backup(command string, args ...string) (startTime time.Time, duration float64, status int, err error) {
	startTime = time.Now()
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()

	defer func() {
		duration = float64(time.Since(startTime).Seconds())
		if err != nil {
			status = 1
			log.Printf("Error:\n", err, string(output))
		} else {
			status = 0
			log.Printf("Command Output: %s\n", string(output))
		}
	}()
	return
}
