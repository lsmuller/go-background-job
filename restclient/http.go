package restclient

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

type RestClient struct {
	logger logrus.FieldLogger
}

func NewRestClient(
	logger logrus.FieldLogger,
) *RestClient {
	logger = logger.WithFields(logrus.Fields{
		"source": "worker",
	})
	return &RestClient{
		logger: logger,
	}
}

func (RestClient) Get(url string) (*http.Response, error) {
	return http.Get(url)
}
