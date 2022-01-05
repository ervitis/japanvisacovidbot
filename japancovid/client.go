package japancovid

import (
	"github.com/go-resty/resty/v2"
	"net/http"
	"time"
)

type (
	restClient struct {
		client *resty.Client
	}

	IRestClient interface {
		R() *resty.Request
	}
)

func NewRestClient() IRestClient {
	c := resty.New()
	c.
		SetHeaders(map[string]string{
			"Content-Type": "application/json",
		}).
		SetRetryCount(3).
		AddRetryCondition(func(response *resty.Response, err error) bool {
			return response.StatusCode() >= http.StatusInternalServerError
		}).
		SetRetryWaitTime(3 * time.Second).
		SetRetryMaxWaitTime(10 * time.Second)

	return &restClient{client: c}
}

func (c *restClient) R() *resty.Request {
	return c.client.R()
}
