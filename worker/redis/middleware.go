package redis

import (
	"github.com/sirupsen/logrus"
	redisWorker "github.com/topfreegames/go-workers"
)

type loggingMiddleware struct {
	logger logrus.FieldLogger
	stack  string
}

func newLoggingMiddleware(
	logger logrus.FieldLogger,
	stack string,
) *loggingMiddleware {
	return &loggingMiddleware{
		logger: logger,
		stack:  stack,
	}
}

func (m *loggingMiddleware) Call(
	queue string,
	message *redisWorker.Msg,
	next func() bool,
) (acknowledge bool) {
	logger := m.logger.WithField("queue", queue)
	logger.Debug("executing task")
	acknowledge = next()
	return
}
