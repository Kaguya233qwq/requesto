package middleware

import (
	"log"
	"time"

	"github.com/Kaguya233qwq/requesto"
)

// RetryPolicy defines the strategy for retrying a failed request.
type RetryPolicy struct {
	RetryCount   int
	RetryBackoff time.Duration
	RetryIf      func(resp *requesto.Response, err error) bool
}

// NewRetrier creates a new retry middleware based on the provided policy.
func NewRetrier(policy RetryPolicy) requesto.Middleware {
	if policy.RetryIf == nil {
		// Default retry condition is to retry on any error.
		policy.RetryIf = func(resp *requesto.Response, err error) bool {
			return err != nil
		}
	}
	if policy.RetryCount <= 0 {
		policy.RetryCount = 3
	}
	if policy.RetryBackoff <= 0 {
		policy.RetryBackoff = 1 * time.Second
	}

	return func(req *requesto.Request, next requesto.Next) (*requesto.Response, error) {
		var resp *requesto.Response
		var err error

		// The loop runs for the initial attempt + RetryCount retries.
		for i := 0; i < policy.RetryCount+1; i++ {
			resp, err = next(req)

			// If the retry condition is not met, return the result immediately.
			if !policy.RetryIf(resp, err) {
				return resp, err
			}

			// If this was the final attempt, break the loop and return the last result.
			if i == policy.RetryCount {
				break
			}

			log.Printf("[Retry] Request failed (%v), retrying in %v...", err, policy.RetryBackoff)
			time.Sleep(policy.RetryBackoff)
		}

		return resp, err
	}
}
