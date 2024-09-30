package file_browser_upload

import "time"

const (
	minimumTimeoutSecond = 3
	minimumRetries       = 3
	defaultConcurrency   = 10
)

type HttpLinkSpeedResult struct {
	URL              string
	BestResponseTime time.Duration
	Error            error
	RetryCount       int // record retries
}

type HttpLinkSpeed struct {
	HttpLinkSpeedFunc HttpLinkSpeedFunc `json:"-"`

	concurrency   uint
	timeoutSecond uint
	maxRetries    uint
}

type HttpLinkSpeedFunc interface {
	SetConcurrency(concurrency uint)

	BestLinkIgnoreRetry(urls []string) (string, error)

	DoUrls(urls []string) ([]HttpLinkSpeedResult, error)
}

// NewLinkSpeed create a new link speed
// timeoutSecond: timeout for each request in second (minimum 5s)
// maxRetries: maxRetries for each request (minimum 3)
func NewLinkSpeed(timeoutSecond, retries uint) *HttpLinkSpeed {
	if timeoutSecond < minimumTimeoutSecond {
		timeoutSecond = minimumTimeoutSecond
	}
	if retries < minimumRetries {
		retries = minimumRetries
	}
	return &HttpLinkSpeed{
		concurrency:   defaultConcurrency,
		timeoutSecond: timeoutSecond,
		maxRetries:    retries,
	}
}
