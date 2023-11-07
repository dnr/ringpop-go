package shared

import (
	"time"

	"golang.org/x/net/context"
)

// var retryOptions = &tchannel.RetryOptions{
// 	RetryOn: tchannel.RetryNever,
// }

// NewTChannelContext creates a new TChannel context with default options
// suitable for use in Ringpop.
func NewTChannelContext(timeout time.Duration) (ContextWithHeaders, context.CancelFunc) {
	panic("unimpl")
	// return tchannel.NewContextBuilder(timeout).
	// 	DisableTracing().
	// 	SetRetryOptions(retryOptions).
	// 	Build()
}

// WrapContext wraps an existing context.Context into a ContextWithHeaders.
// If the underlying context has headers, they are preserved.
func WrapContext(ctx context.Context) ContextWithHeaders {
	panic("unimpl")
}
