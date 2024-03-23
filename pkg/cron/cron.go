package cron

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	expr "github.com/gitploy-io/cronexpr"
)

const SingleRun = "single_run"

type Cron interface {
	AddJob(name, tm string, cmd Cmd) error
	Run(ctx context.Context) error
}

type Cmd func(ctx context.Context)

func NewCron() Cron {
	return &cron{
		wg: new(sync.WaitGroup),
	}
}

type cron struct {
	jobs []job
	wg   *sync.WaitGroup
}

type job struct {
	name string
	cmd  Cmd
	shed Schedule
}

type Schedule interface {
	Next(t time.Time) time.Time
}

func (c *cron) AddJob(name, tm string, cmd Cmd) error {
	schedule, err := parseTime(tm)
	if err != nil {
		return fmt.Errorf("error parsing time: %w", err)
	}
	c.jobs = append(c.jobs, job{
		name: name,
		cmd:  cmd,
		shed: schedule,
	})
	return nil
}

func parseTime(tm string) (schedule Schedule, err error) {
	if tm == SingleRun {
		return new(singleRunner), nil
	}
	if schedule, err = expr.Parse(tm); err != nil {
		return nil, fmt.Errorf("error parsing cron time: %w", err)
	}
	return
}

func (c *cron) Run(ctx context.Context) error {
	logger := log.New(os.Stdout, "", 0)
	logger.Println("cron started")
	for _, j := range c.jobs {
		j := j
		c.wg.Add(1)
		go func() {
			defer c.wg.Done()
			now := time.Now()
			nextExec := j.shed.Next(now)
			logger.Printf("job: %s next exec at: %s\n", j.name, nextExec.Format(time.RFC3339))
			timer := time.NewTimer(nextExec.Sub(now))
			for {
				select {
				case <-timer.C:
					now := time.Now()
					logger.Printf("started job: %s\n", j.name)
					j.cmd(ctx)
					nextExec := j.shed.Next(now)
					logger.Printf("finished job: %s, next exec at: %s\n", j.name, nextExec.Format(time.RFC3339))
					timer.Reset(nextExec.Sub(now))
				case <-ctx.Done(): // exit job runner
					return
				}
			}
		}()
	}
	<-ctx.Done() // exit cron
	logger.Println("cron exiting")
	c.wg.Wait()
	logger.Println("cron all jobs stopped")
	return nil
}

// singleRunner runs job now and next is in 100 years
// expected to use for daemon services
type singleRunner struct {
	done bool
}

func (sr *singleRunner) Next(t time.Time) time.Time {
	if !sr.done {
		sr.done = true
		return t
	}
	return t.AddDate(100, 0, 0)
}
