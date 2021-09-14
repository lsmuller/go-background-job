package redis

import (
	"github.com/lsmuller/go-background-job/config"
	redisWorker "github.com/topfreegames/go-workers"
)

type QueueConfig struct {
	Retry    bool
	RetryMax int
	Exp      int
	MinDelay int
	MaxDelay int
	MaxRand  int
}

func NewQueueConfig(config *config.WorkerConfig) *QueueConfig {
	return &QueueConfig{
		Retry:    config.Retry.Enabled,
		RetryMax: config.Retry.RetryMax,
		Exp:      config.Retry.Exp,
		MinDelay: config.Retry.MinDelay,
		MaxDelay: config.Retry.MaxDelay,
		MaxRand:  config.Retry.MaxRand,
	}
}

func (r *QueueConfig) ToRedisEnqueueOptions() *redisWorker.EnqueueOptions {
	return &redisWorker.EnqueueOptions{
		Retry:    r.Retry,
		RetryMax: r.RetryMax,
		RetryOptions: redisWorker.RetryOptions{
			Exp:      r.Exp,
			MinDelay: r.MinDelay,
			MaxDelay: r.MaxDelay,
			MaxRand:  r.MaxRand,
		},
	}
}
