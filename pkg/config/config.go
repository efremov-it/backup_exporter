package config

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	Host, Port     string
	ProjectName    string
	InstanceName   string
	BackupType     string
	ConfigFile     string
	BackupStorage  string
	CronTime       string
	DeleteRetain   string
	DeleteCronTime string
}

func getHostname() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}
	return hostname, nil
}

func ParseFlags() (*Config, error) {
	config := &Config{}
	// Define flags
	flag.StringVar(&config.Host, "host", "127.0.0.1", "Backup server address (default: \"127.0.0.1\")")
	flag.StringVar(&config.Port, "port", "9023", "Backup server port (default: \"9023\")")
	flag.StringVar(&config.ProjectName, "project", "", "Project name")
	flag.StringVar(&config.InstanceName, "instance", "", "Instance name (default is the hostname)")
	flag.StringVar(&config.BackupType, "backup_type", "", "Backup type (psql|mysql|mariadb|mongodb|clickhouse)")
	flag.StringVar(&config.ConfigFile, "config_file", "", "Path to config file (default is different for each backup tools)")
	flag.StringVar(&config.BackupStorage, "backup_storage", "", "When uploading backups to storage, the user should pass the Postgres data directory as an argument (default: $PGDATA)")
	flag.StringVar(&config.CronTime, "backup_cron", "0 2 * * *", "How often you should create your backup. Format crontab --backup_cron \"* * * * *\" (default: \"0 2 * * *\" At 02:00)")
	flag.StringVar(&config.DeleteRetain, "delete_retain", "", "when set keep 5(or) full backups and all deltas of them (default: not set. example 5)")
	flag.StringVar(&config.DeleteCronTime, "delete_cron", "0 3* * 6", "How often you should retain your backups (default: \"0 3* * 6\" At 03:00 on Saturday)")

	flag.Parse()

	errUsage := fmt.Errorf("usage: main --project projectName --backup_type <psql|mysql|mariadb|mongodb|clickhouse> --cron \"* * * * *\"")

	fmt.Println("DeleteRetain \n DeleteCronTime", config.DeleteRetain, config.DeleteCronTime)
	// Set default values

	if config.InstanceName == "" {
		hostname, err := getHostname()
		if err != nil {
			return nil, err
		}
		config.InstanceName = hostname
	}

	if config.BackupStorage == "" && config.BackupType == "psql" {
		pgdata := os.Getenv("PGDATA")
		if pgdata == "" {
			return nil, fmt.Errorf("error: $PGDATA environment variable is not set")
		}
		config.BackupStorage = pgdata
	}

	checkBackupType := true
	for _, db := range []string{"psql", "mysql", "mariadb", "mongodb", "clickhouse"} {
		if db == config.BackupType {
			checkBackupType = false
		}
	}

	// Validate required flags
	if config.ProjectName == "" || checkBackupType {
		return nil, errUsage
	}

	return config, nil
}
