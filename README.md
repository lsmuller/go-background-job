This projects enables quickly creation of periodic background jobs.


Either way, worker stands for background jobs executed by the server. These processes run behind the scenes and without the player's intervention. Typical tasks for these processes include logging, system monitoring, scheduling, and user notification.

At this point of the lifecycle of the game, for whatever we code we should think about scaling. Millions of players use this game and a huge amount of data flows within what we code. With that in mind, we have chosen to work with the package [go-workers](https://github.com/topfreegames/go-workers). It uses queues in [Redis](https://redis.io/) to trigger workers, which doesn't encumber our server with such memory and processing.

### Step by step to create a worker

**1. Define your worked and add a constructor.**  
The structure should consider all the relevant attributes and structures you will need within the worker for it to be functional (DALs, BLLs, services, ...).

```
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
```

**2. Create the Start and Stop methods**   
`Start`: receives the main game configuration and uses it to set all the necessary information for the job to execute. Also, registers the jobs the worker is responsible for.  
`Stop`: Updates isRunning to false, which exits the worker main loop.
```
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
```
Make sure you added all necessary fields to the config. `enabled`, `concurrency` and `period` are default. Use `metadata` to any further config necessary.
```
    bitcoinPriceIndexFetcher:
      enabled: true
      concurrency: 20
      period: "1m"
      metadata:
        apiUrl: https://api.coindesk.com/v1/bpi/currentprice.json
```

**3. Create the job to do the necessary processing**  
In this example there are no arguments being expected by the method in order to execute. It is possible to pass arguments from the scheduler using the jobArg parameter. There are examples on how to unmarshall it in existing workers of the project.  
This is where you should interact with any other layers of the application and perform the behavior expected for the job.
```
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
```

**4. Register your created job**
- Create a constant to represent the queue the worker will use to trigger the job. File: `constants/worker`
```
const (
    ...
	BitcoinPriceIndexFetcher = "bitcoinPriceIndexFetcher"
)
```
- Use the register method to assign the job you created with the queue. By doing so, whenever there's a new message at the queue, it will be forwarded to the assign method that will know how to parse it and process it.
```
func (b *BitcoinPriceIndexFetcherWorker) register(concurrency int) {
	b.worker.Register(constants.BitcoinPriceIndexFetcher, concurrency, b.fetchBitcoinPriceIndex)
}
```

**5. Enqueue to allow async processing**  
With the queue created and the job registered, you just need to enqueue the job and the lib and redis will do the rest.
```
func (b *BitcoinPriceIndexFetcherWorker) enqueue() error {
	if err := b.worker.Enqueue(constants.BitcoinPriceIndexFetcher, nil, b.queueConfig.ToRedisEnqueueOptions()); err != nil {
		return err
	}
	return nil
}
```
**6. Create periodical execution**
Create a scheduler to trigger the worker
```
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
```

7. Start the worker in `main.go` (start of the application).
```
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
```

Testing
```bash
#Setup project dependencies
make deps

#Run server on worker mode
make run

#CTRL C / COMMAND C to stop
```

Output
```
{"file":"~/go/src/github.com/lsmuller/go-background-jobs/worker/bitcoin_price_fetcher.go:109","func":"github.com/lsmuller/go-background-jobs/worker.(*BitcoinPriceIndexFetcherWorker).fetchBitcoinPriceIndex","level":"info","method":"fetchBitcoinPriceIndex","msg":"On Sep 14, 2021 01:13:00 UTC Bitcoin value is worth 45,209.7524 USD","time":"2021-09-13T22:13:48-03:00"}
{"file":"~/go/src/github.com/lsmuller/go-background-jobs/worker/bitcoin_price_fetcher.go:110","func":"github.com/lsmuller/go-background-jobs/worker.(*BitcoinPriceIndexFetcherWorker).fetchBitcoinPriceIndex","level":"warning","method":"fetchBitcoinPriceIndex","msg":"This data was produced from the CoinDesk Bitcoin Price Index (USD). Non-USD currency data converted using hourly conversion rate from openexchangerates.org","time":"2021-09-13T22:13:48-03:00"}
```