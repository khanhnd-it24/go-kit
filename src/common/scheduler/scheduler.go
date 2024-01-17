package schedulerprovider

import (
	"context"
	"github.com/go-co-op/gocron/v2"
)

type SchedulerProvider struct {
	Scheduler gocron.Scheduler
}

func NewSchedulerProvider() (*SchedulerProvider, error) {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}
	return &SchedulerProvider{
		Scheduler: scheduler,
	}, nil
}

func (s *SchedulerProvider) Start(ctx context.Context) error {
	s.Scheduler.Start()
	return nil
}

func (s *SchedulerProvider) Stop(ctx context.Context) error {
	if err := s.Scheduler.Shutdown(); err != nil {
		return err
	}
	return nil
}

func (s *SchedulerProvider) NewJob(crontab string, task gocron.Task) (string, error) {
	j, err := s.Scheduler.NewJob(gocron.CronJob(crontab, false), task)
	if err != nil {
		return "", err
	}
	return j.ID().String(), nil
}
