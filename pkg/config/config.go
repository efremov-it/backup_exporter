package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Host, Port     string
	ProjectName    string
	InstanceName   string
	BackupType     string
	ConfigFile     string
	BackupStorage  string
	CronTime       string
	DeleteEnable   bool
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
	flag.StringVar(&config.Host, "host", "127.0.0.1", "Backup server address")
	flag.StringVar(&config.Port, "port", "9023", "Backup server port")
	flag.StringVar(&config.ProjectName, "project", "", "Project name (Necessarily)")
	flag.StringVar(&config.InstanceName, "instance", "", "Instance name (default is the hostname)")
	flag.StringVar(&config.BackupType, "backup_type", "", "Backup type (psql|mysql|mariadb|mongodb|clickhouse) (Necessarily)")
	flag.StringVar(&config.ConfigFile, "config_file", "", "Path to config file (default is different for each backup tools)")
	flag.StringVar(&config.BackupStorage, "backup_storage", "", "When uploading backups to storage, the user should pass the Postgres data directory as an argument (default: $PGDATA)")
	flag.StringVar(&config.CronTime, "backup_cron", "0 2 * * *", "How often you should create your backup. Format crontab --backup_cron \"* * * * *\"")
	flag.StringVar(&config.DeleteRetain, "delete_retain", "", "when set keep 5(or) full backups and all deltas of them (default: not set. example 5 int type)")
	flag.StringVar(&config.DeleteCronTime, "delete_cron", "0 3* * 6", "How often you should retain your backups")

	flag.Parse()

	errUsage := fmt.Errorf("\n-----------------------------\nusage: main --project projectName --backup_type <psql|mysql|mariadb|mongodb|clickhouse> --backup_cron \"* * * * *\"\nSupport only for postgresql (--delete_cron \"*/1 * * * *\" --delete_retain 5)")
	config.DeleteEnable = false

	if config.ProjectName == "" {
		return nil, errUsage
	}

	switch config.BackupType {
	case "psql":
		if config.BackupStorage == "" {
			config.BackupStorage = os.Getenv("PGDATA")
		}
		if config.BackupStorage == "" {
			return nil, fmt.Errorf("error: $PGDATA environment variable is not set")
		}
		// enable backup Retain
		if config.DeleteRetain != "" {
			// if delete_retain contain string
			_, err := strconv.Atoi(config.DeleteRetain)
			if err != nil {
				return nil, errUsage
			}
			config.DeleteEnable = true
		} else {
			fmt.Print("For enable deleting old backups you need to use flags:\n --delete_cron \"*/1 * * * *\" --delete_retain 5\n")
		}
		fmt.Printf("PostgreSQL backup selected\n")
	case "mysql":
		fmt.Println("MySQL backup selected")
	case "mariadb":
		fmt.Println("MariaDB backup selected")
	case "mongodb":
		fmt.Println("MongoDB backup selected")
	case "clickhouse":
		fmt.Println("ClickHouse backup selected")
	default:
		fmt.Println("Invalid or no backup type specified")
		return nil, errUsage
	}

	// Set up default values
	if config.InstanceName == "" {
		hostname, err := getHostname()
		if err != nil {
			return nil, err
		}
		config.InstanceName = hostname
	}

	return config, nil
}
