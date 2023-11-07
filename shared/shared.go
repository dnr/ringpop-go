package shared

import (
	"time"

	"context"
)

type (
	contextKeyHeadersType struct{}
)

var contextKeyHeaders = contextKeyHeadersType{}

// FIXME
// var retryOptions = &tchannel.RetryOptions{
// 	RetryOn: tchannel.RetryNever,
// }

// NewTChannelContext creates a new TChannel context with default options
// suitable for use in Ringpop.
func NewTChannelContext(timeout time.Duration) (ContextWithHeaders, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	return WrapContext(ctx), cancel
}

type headerCtx struct {
	context.Context
}

// headersContainer stores the headers, and is itself stored in the context under `contextKeyHeaders`
type headersContainer struct {
	reqHeaders  map[string]string
	respHeaders map[string]string
}

func (c headerCtx) headers() *headersContainer {
	if h, ok := c.Value(contextKeyHeaders).(*headersContainer); ok {
		return h
	}
	return nil
}

// Headers gets application headers out of the context.
func (c headerCtx) Headers() map[string]string {
	if h := c.headers(); h != nil {
		return h.reqHeaders
	}
	return nil
}

// ResponseHeaders returns the response headers.
func (c headerCtx) ResponseHeaders() map[string]string {
	if h := c.headers(); h != nil {
		return h.respHeaders
	}
	return nil
}

// SetResponseHeaders sets the response headers.
func (c headerCtx) SetResponseHeaders(headers map[string]string) {
	if h := c.headers(); h != nil {
		h.respHeaders = headers
		return
	}
	panic("SetResponseHeaders called on ContextWithHeaders not created via WrapWithHeaders")
}

// Child creates a child context with a separate container for headers.
func (c headerCtx) Child() ContextWithHeaders {
	var headersCopy headersContainer
	if h := c.headers(); h != nil {
		headersCopy = *h
	}

	return WrapContext(context.WithValue(c.Context, contextKeyHeaders, &headersCopy))
}

// WrapContext wraps an existing context.Context into a ContextWithHeaders.
// If the underlying context has headers, they are preserved.
func WrapContext(ctx context.Context) ContextWithHeaders {
	hctx := headerCtx{Context: ctx}
	if h := hctx.headers(); h != nil {
		return hctx
	}

	// If there is no header container, we should create an empty one.
	return WrapWithHeaders(ctx, nil)
}

// WrapWithHeaders returns a Context that can be used to make a call with request headers.
// If the parent `ctx` is already an instance of ContextWithHeaders, its existing headers
// will be ignored. In order to merge new headers with parent headers, use ContextBuilder.
func WrapWithHeaders(ctx context.Context, headers map[string]string) ContextWithHeaders {
	h := &headersContainer{
		reqHeaders: headers,
	}
	newCtx := context.WithValue(ctx, contextKeyHeaders, h)
	return headerCtx{Context: newCtx}
}
