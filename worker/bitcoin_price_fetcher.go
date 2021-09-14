package worker

import (
	"encoding/json"
	"github.com/lsmuller/go-background-job/config"
	"github.com/lsmuller/go-background-job/constants"
	"github.com/lsmuller/go-background-job/worker/redis"
	"github.com/sirupsen/logrus"
	redisWorker "github.com/topfreegames/go-workers"
	"io/ioutil"
	"net/http"
	"time"
)

type BitcoinPriceIndexFetcherWorker struct {
	logger      logrus.FieldLogger
	queueConfig *redis.QueueConfig
	worker      redis.WorkerInterface
	apiUrl      string
	isRunning   bool
}

func NewBitcoinPriceIndexFetcherWorker(
	queueConfig *redis.QueueConfig,
	worker redis.WorkerInterface,
	logger logrus.FieldLogger,
) *BitcoinPriceIndexFetcherWorker {
	return &BitcoinPriceIndexFetcherWorker{
		logger:      logger,
		queueConfig: queueConfig,
		worker:      worker,
	}
}

func (b *BitcoinPriceIndexFetcherWorker) Start(config *config.WorkerConfig) {
	job := config.Jobs.BitcoinPriceIndexFetcher
	metadata := job.Metadata
	b.apiUrl = metadata["apiUrl"]

	b.register(job.Concurrency)
	b.isRunning = job.Enabled

	go b.scheduler(job.Period)
}

func (b *BitcoinPriceIndexFetcherWorker) Stop() {
	b.isRunning = false
}

func (b *BitcoinPriceIndexFetcherWorker) register(concurrency int) {
	b.worker.Register(constants.BitcoinPriceIndexFetcher, concurrency, b.fetchBitcoinPriceIndex)
}

type BitcoinPriceIndexData struct {
	Time struct {
		Updated    string    `json:"updated"`
		UpdatedISO time.Time `json:"updatedISO"`
		Updateduk  string    `json:"updateduk"`
	} `json:"time"`
	Disclaimer string `json:"disclaimer"`
	ChartName  string `json:"chartName"`
	Bpi        struct {
		USD struct {
			Code        string  `json:"code"`
			Symbol      string  `json:"symbol"`
			Rate        string  `json:"rate"`
			Description string  `json:"description"`
			RateFloat   float64 `json:"rate_float"`
		} `json:"USD"`
		GBP struct {
			Code        string  `json:"code"`
			Symbol      string  `json:"symbol"`
			Rate        string  `json:"rate"`
			Description string  `json:"description"`
			RateFloat   float64 `json:"rate_float"`
		} `json:"GBP"`
		EUR struct {
			Code        string  `json:"code"`
			Symbol      string  `json:"symbol"`
			Rate        string  `json:"rate"`
			Description string  `json:"description"`
			RateFloat   float64 `json:"rate_float"`
		} `json:"EUR"`
	} `json:"bpi"`
}

// addEventRankRewards Enqueues the creation of leaderboard rewards for the given player group
func (b *BitcoinPriceIndexFetcherWorker) fetchBitcoinPriceIndex(jobArg *redisWorker.Msg) {
	logger := b.logger.WithFields(logrus.Fields{
		"method": "fetchBitcoinPriceIndex",
	})

	resp, err := http.Get(b.apiUrl)
	if err != nil {
		logger.Error(err)
	}

	defer resp.Body.Close()

	jsonDataFromHttp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var data *BitcoinPriceIndexData

	err = json.Unmarshal(jsonDataFromHttp, &data)

	logger.Info("On " + data.Time.Updated + " Bitcoin value is worth " + data.Bpi.USD.Rate + " " + data.Bpi.USD.Code)
	logger.Warn(data.Disclaimer)

}

// CreateNewEvents Enqueues the creation of new events
func (b *BitcoinPriceIndexFetcherWorker) enqueue() error {
	if err := b.worker.Enqueue(constants.BitcoinPriceIndexFetcher, nil, b.queueConfig.ToRedisEnqueueOptions()); err != nil {
		return err
	}
	return nil
}

// PeriodicEventScheduler periodically creates tasks to create events or calculate rewards
func (b *BitcoinPriceIndexFetcherWorker) scheduler(period time.Duration) {
	logger := b.logger.WithFields(logrus.Fields{
		"method": "BitcoinPriceIndexFetcherScheduler",
	})
	for b.isRunning {
		time.Sleep(period)
		logger.Debugf("[BitcoinPriceIndexFetcherWorker] Running BitcoinPriceIndexFetcherScheduler")

		if err := b.enqueue(); err != nil {
			logger.WithError(err).Error("[BitcoinPriceIndexFetcherWorker] Failed to fetch bitcoin prices")
			continue
		}
	}
}
