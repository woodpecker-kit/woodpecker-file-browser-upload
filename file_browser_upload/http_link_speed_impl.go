package file_browser_upload

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// SetConcurrency set concurrency for link speed
// minimum concurrency is defaultConcurrency = 10
func (l *HttpLinkSpeed) SetConcurrency(concurrency uint) {
	if concurrency > defaultConcurrency {
		l.concurrency = concurrency
	}
}

func (l *HttpLinkSpeed) BestLinkIgnoreRetry(urls []string) (string, error) {
	if len(urls) == 0 {
		return "", fmt.Errorf("urls is empty")
	}

	results, errDoUrls := l.DoUrls(urls)
	if errDoUrls != nil {
		return "", errDoUrls
	}

	bestURL := ""
	bestResponseTime := time.Duration(0)
	for _, result := range results {
		if result.BestResponseTime == 0 {
			continue
		}
		if bestResponseTime == 0 || result.BestResponseTime < bestResponseTime {
			bestResponseTime = result.BestResponseTime
			bestURL = result.URL
		}
	}

	if bestURL == "" {
		return "", fmt.Errorf("not found best url")
	}

	return bestURL, nil
}

func (l *HttpLinkSpeed) DoUrls(urls []string) ([]HttpLinkSpeedResult, error) {
	if len(urls) == 0 {
		return nil, fmt.Errorf("urls is empty")
	}

	var wg sync.WaitGroup
	results := make(chan HttpLinkSpeedResult, len(urls))

	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()

			var result HttpLinkSpeedResult
			var ctx context.Context
			var cancel context.CancelFunc

			maxRetry := int(l.maxRetries)
			timeoutSecond := time.Duration(l.timeoutSecond) * time.Second
			bestResponseTime := time.Duration(0)

			for retries := 0; retries <= maxRetry; retries++ {
				result.URL = u

				ctx, cancel = context.WithTimeout(context.Background(), timeoutSecond)
				req, errNewRequest := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
				if errNewRequest != nil {
					result.Error = errNewRequest
					result.RetryCount = result.RetryCount + 1
					cancel() // 取消当前请求的上下文，释放资源
					continue
				}

				req.Header.Add("Cache-Control", "no-cache")
				start := time.Now()
				resp, errDoRequest := http.DefaultClient.Do(req)
				if errDoRequest != nil {
					result.Error = errDoRequest
					result.RetryCount = result.RetryCount + 1
					cancel() // 取消当前请求的上下文，释放资源
					continue
				}
				defer resp.Body.Close()

				durationTIme := time.Since(start)

				if bestResponseTime == 0 || durationTIme < bestResponseTime {
					bestResponseTime = durationTIme
				}

				result.BestResponseTime = bestResponseTime
				result.Error = nil
				cancel() // 取消当前请求的上下文，释放资源
			}

			results <- result
		}(url)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var allResults []HttpLinkSpeedResult
	for result := range results {
		allResults = append(allResults, result)
	}

	return allResults, nil
}
