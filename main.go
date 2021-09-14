package main

import (
	"github.com/lsmuller/go-background-job/config"
	"github.com/lsmuller/go-background-job/logging"
	serverWorker "github.com/lsmuller/go-background-job/worker"
	redisWorker "github.com/lsmuller/go-background-job/worker/redis"
)

func main() {
	cfg := config.GetConfig()
	logger := logging.CreateLogger(cfg.App.Logging)

	worker, err := redisWorker.NewWorker(cfg, logger)
	if err != nil {
		logger.WithError(err).Fatalf("failed to create worker")
	}
	worker.Setup()

	queueConfig := redisWorker.NewQueueConfig(cfg.Worker)

	bitcoinPriceIndexFetcher := serverWorker.NewBitcoinPriceIndexFetcherWorker(queueConfig, worker, logger)
	bitcoinPriceIndexFetcher.Start(cfg.Worker)

	worker.Start()
}
