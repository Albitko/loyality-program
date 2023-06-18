package utils

import (
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

var RestyClient *resty.Client

func InitializeRestyClient() {
	RestyClient = resty.New()

	RestyClient.
		AddRetryCondition(
			func(r *resty.Response, err error) bool {
				return r.StatusCode() == http.StatusTooManyRequests
			},
		).
		SetRetryMaxWaitTime(120 * time.Second)

}
