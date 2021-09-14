package redis

import (
	"github.com/lsmuller/go-background-job/config"
	"github.com/sirupsen/logrus"
	redisWorker "github.com/topfreegames/go-workers"
	"os"
)

type WorkerInterface interface {
	Register(queue string, concurrency int, job func(jobArg *redisWorker.Msg))
	Enqueue(queue string, args interface{}, opts *redisWorker.EnqueueOptions) error
}

// Worker implements WorkerInterface
type Worker struct {
	cfg    *config.MainConfig
	logger logrus.FieldLogger
}

// NewWorker creates a Worker instance.
func NewWorker(
	cfg *config.MainConfig,
	logger logrus.FieldLogger,
) (*Worker, error) {
	logger = logger.WithFields(logrus.Fields{
		"source": "worker",
	})
	return &Worker{
		cfg:    cfg,
		logger: logger,
	}, nil
}

//Setup Initializes the basic config required to communicate with Redis
func (w *Worker) Setup() {
	hostname, _ := os.Hostname()

	redisWorker.Configure(map[string]string{
		"server":    w.cfg.Worker.Redis.Url,
		"pool":      w.cfg.Worker.Redis.Pool,
		"namespace": w.cfg.App.Stack,
		"process":   hostname,
	})
	redisWorker.SetLogger(w.logger)
	redisWorker.Middleware.Append(newLoggingMiddleware(w.logger, w.cfg.App.Stack))

	go redisWorker.StatsServer(8080)
}

// Register Links the receiving job to the receiving queue
func (w *Worker) Register(queue string, concurrency int, job func(jobArg *redisWorker.Msg)) {
	redisWorker.Process(queue, job, concurrency)
}

// Enqueue Enqueues the receiving args in the receiving queue to be executed by the middleware
func (w *Worker) Enqueue(queue string, args interface{}, opts *redisWorker.EnqueueOptions) error {
	if _, err := redisWorker.EnqueueWithOptions(queue, "", args, *opts); err != nil {
		return err
	}
	return nil
}

func (w *Worker) Start() {
	redisWorker.Run()
}
