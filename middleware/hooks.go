package middleware

import "github.com/Kaguya233qwq/requesto"

// BeforeRequestFunc defines the signature for a hook function that runs before a request is sent.
type BeforeRequestFunc func(req *requesto.Request) error

// AfterResponseFunc defines the signature for a hook function that runs after a response is received.
type AfterResponseFunc func(resp *requesto.Response, err error)

// hookConfig is a private struct used to store the configuration for the NewHook middleware.
type hookConfig struct {
	before []BeforeRequestFunc
	after  []AfterResponseFunc
}

// HookOption is a function type used to configure the NewHook middleware builder.
type HookOption func(*hookConfig)

// WithBeforeRequest adds one or more hooks to be executed before the request is sent.
func WithBeforeRequest(hooks ...BeforeRequestFunc) HookOption {
	return func(c *hookConfig) {
		c.before = append(c.before, hooks...)
	}
}

// WithAfterResponse adds one or more hooks to be executed after the response is received.
func WithAfterResponse(hooks ...AfterResponseFunc) HookOption {
	return func(c *hookConfig) {
		c.after = append(c.after, hooks...)
	}
}

// NewHook is a middleware builder that constructs a full requesto.Middleware
// from simple hook functions (before request and after response).
func NewHook(opts ...HookOption) requesto.Middleware {
	// Create a default configuration.
	config := &hookConfig{
		before: make([]BeforeRequestFunc, 0),
		after:  make([]AfterResponseFunc, 0),
	}

	// Apply all user-provided options.
	for _, opt := range opts {
		opt(config)
	}

	return func(req *requesto.Request, next requesto.Next) (*requesto.Response, error) {
		// Execute before-request hooks.
		for _, hook := range config.before {
			if err := hook(req); err != nil {
				// If any before-hook returns an error, abort the request chain immediately.
				return nil, err
			}
		}

		// Proceed to the next middleware or the actual request.
		resp, err := next(req)

		// Execute after-response hooks in reverse order.
		for i := len(config.after) - 1; i >= 0; i-- {
			config.after[i](resp, err)
		}

		return resp, err
	}
}
