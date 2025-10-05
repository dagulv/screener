package cron

import "github.com/go-co-op/gocron/v2"

func New() (s gocron.Scheduler, err error) {
	s, err = gocron.NewScheduler()
	if err != nil {
		return
	}
	return
}
