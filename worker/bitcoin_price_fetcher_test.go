package worker_test

import (
	"github.com/golang/mock/gomock"
	"github.com/lsmuller/go-background-job/config"
	"github.com/lsmuller/go-background-job/constants"
	"github.com/lsmuller/go-background-job/mocks"
	"github.com/lsmuller/go-background-job/test"
	"github.com/lsmuller/go-background-job/worker"
	"github.com/lsmuller/go-background-job/worker/redis"
	"github.com/sirupsen/logrus"
	testLogger "github.com/sirupsen/logrus/hooks/test"
	redisWorker "github.com/topfreegames/go-workers"
	"net/http"
	"testing"
	"time"
)

func TestBitcoinPriceIndexFetcherWorker(t *testing.T) {
	t.Parallel()

	var (
		cfg            *config.MainConfig
		logger         *logrus.Logger
		queueConfig    *redis.QueueConfig
		mockWorker     *mocks.MockWorkerInterface
		mockHttpClient *mocks.MockHTTPClient
	)

	setup := func(ctrl *gomock.Controller) {
		cfg = test.GetConfig()
		logger, _ = testLogger.NewNullLogger()
		queueConfig = redis.NewQueueConfig(cfg.Worker)
		mockWorker = mocks.NewMockWorkerInterface(ctrl)
		mockHttpClient = mocks.NewMockHTTPClient(ctrl)
	}

	t.Run("worker_disabled", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		setup(ctrl)

		cfg.Worker.Jobs.BitcoinPriceIndexFetcher.Enabled = false

		mockWorker.EXPECT().Register(constants.BitcoinPriceIndexFetcher, gomock.Any(), gomock.Any())

		bitcoinPriceIndexFetcher := worker.NewBitcoinPriceIndexFetcherWorker(queueConfig, mockWorker, logger, mockHttpClient)
		bitcoinPriceIndexFetcher.Start(cfg.Worker)
		bitcoinPriceIndexFetcher.Stop()
	})

	t.Run("success_enqueue", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		setup(ctrl)

		cfg.Worker.Jobs.BitcoinPriceIndexFetcher.Period, _ = time.ParseDuration("0s")

		mockWorker.EXPECT().Register(constants.BitcoinPriceIndexFetcher, gomock.Any(), gomock.Any()).AnyTimes()
		mockWorker.EXPECT().Enqueue(constants.BitcoinPriceIndexFetcher, gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

		bitcoinPriceIndexFetcher := worker.NewBitcoinPriceIndexFetcherWorker(queueConfig, mockWorker, logger, mockHttpClient)
		bitcoinPriceIndexFetcher.Start(cfg.Worker)
		bitcoinPriceIndexFetcher.Stop()
	})

	t.Run("fetch_bitcoin_price_index", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		setup(ctrl)

		mockHttpClient.EXPECT().Get(gomock.Any()).Return(&http.Response{StatusCode: 200}, nil)

		bitcoinPriceIndexFetcher := worker.NewBitcoinPriceIndexFetcherWorker(queueConfig, mockWorker, logger, mockHttpClient)

		bitcoinPriceIndexFetcher.FetchBitcoinPriceIndex(&redisWorker.Msg{})
	})
}
