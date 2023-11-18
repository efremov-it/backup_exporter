package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"os"

	"github.com/efremov-it/backup_exporter/pkg/collector"
	"github.com/efremov-it/backup_exporter/pkg/config"
	"github.com/efremov-it/backup_exporter/pkg/create"
	"github.com/efremov-it/backup_exporter/pkg/cron"
)

// c := NewCron()
// if err := c.AddJob("First", "* * * * *", func(ctx context.Context) {
// }); err != nil {
//  panic(err)
// }
// if err := c.AddJob("Second", "* * * * *", func(ctx context.Context) {
// }); err != nil {
//  panic(err)
// }

func main() {
	config, err := config.ParseFlags()

	if err != nil {
		log.Fatal("Error parsing flags:\n", err)
	}

	c := cron.NewCron()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()
	collectorVars := collector.NewCollector(config)
	prometheus.MustRegister(collectorVars)
	backupService := create.NewBackupService(config, collectorVars)

	if err := c.AddJob("Backup", config.CronTime, backupService.Run); err != nil {
		panic(err)
	}
	// create.NewBackupService(config, collectorVars).Run(ctx)
	// Create a Prometheus registry
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	s := http.Server{Addr: fmt.Sprintf("%s:%s", config.Host, config.Port), Handler: mux}
	if err := c.AddJob("MetricServer", cron.SingleRun, func(ctx context.Context) {
		go func() {
			<-ctx.Done()
			s.Shutdown(context.Background())
		}()
		err := s.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}); err != nil {
		panic(err)
	}

	if err := c.Run(ctx); err != nil {
		panic(err)
	}

}
