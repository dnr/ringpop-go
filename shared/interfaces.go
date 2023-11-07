package shared

import (
	"context"
)

// The TChannel interface defines the dependencies for TChannel in Ringpop.
type TChannel interface {
	GetSubChannel(serviceName string) SubChannel
	PeerInfo() LocalPeerInfo
	ListenAndServe(hostport string) error
	Ping(context.Context, string) error
	Close()
}

// SubChannel represents a TChannel SubChannel as used in Ringpop.
type SubChannel interface {
	Registrar
}

// Handlers is the map from method names to handlers.
type Handlers map[string]interface{}

// Registrar is the base interface for registering handlers on either the base
// Channel or the SubChannel
type Registrar interface {
	// ServiceName returns the service name that this Registrar is for.
	ServiceName() string

	// Register registers the specified methods specified as a map from method name to the
	// JSON handler function. The handler functions should have the following signature:
	// func(context.Context, *ArgType)(*ResType, error)
	Register(funcs Handlers, onError func(context.Context, error)) error

	// // Logger returns the logger for this Registrar.
	// Logger() Logger

	// // StatsReporter returns the stats reporter for this Registrar
	// StatsReporter() StatsReporter

	// // StatsTags returns the tags that should be used.
	// StatsTags() map[string]string

	// Peers returns the peer list for this Registrar.
	Peers() PeerList
}

// A Handler is an object that can be registered with a Channel to process
// incoming calls for a given service and method
type Handler interface {
	// // Handles an incoming call for service
	// Handle(ctx context.Context, call InboundCall)
}

// Format is the arg scheme used for a specific call.
type Format string

// The list of formats supported by tchannel.
const (
	HTTP   Format = "http"
	JSON   Format = "json"
	Raw    Format = "raw"
	Thrift Format = "thrift"
)

// CallOptions are options for a specific call.
type CallOptions struct {
	RawHeaders []byte // FIXME
	Headers    any    // FIXME

	// // Format is arg scheme used for this call, sent in the "as" header.
	// // This header is only set if the Format is set.
	// Format Format

	// // ShardKey determines where this call request belongs, used with ringpop applications.
	// ShardKey string

	// // RequestState stores request state across retry attempts.
	// RequestState *tchannel.RequestState

	// // RoutingKey identifies the destined traffic group. Relays may favor the
	// // routing key over the service name to route the request to a specialized
	// // traffic group.
	// RoutingKey string

	// // RoutingDelegate identifies a traffic group capable of routing a request
	// // to an instance of the intended service.
	// RoutingDelegate string

	// // CallerName defaults to the channel's service name for an outbound call.
	// // Optionally override this field to support transparent proxying when inbound
	// // caller names vary across calls.
	// CallerName string
}

// // ArgReader is the interface for the arg2 and arg3 streams on an
// // OutboundCallResponse and an InboundCall
// type ArgReader io.ReadCloser

// // ArgWriter is the interface for the arg2 and arg3 streams on an OutboundCall
// // and an InboundCallResponse
// type ArgWriter interface {
// 	io.WriteCloser

// 	// Flush flushes the currently written bytes without waiting for the frame
// 	// to be filled.
// 	Flush() error
// }

// An InboundCall is an incoming call from a peer
type InboundCall interface {
	// // ServiceName returns the name of the service being called
	// ServiceName() string
	// // Method returns the method being called
	// Method() []byte
	// // MethodString returns the method being called as a string.
	// MethodString() string
	// // Format the format of the request from the ArgScheme transport header.
	// Format() Format
	// // CallerName returns the caller name from the CallerName transport header.
	// CallerName() string
	// // ShardKey returns the shard key from the ShardKey transport header.
	// ShardKey() string
	// // RoutingKey returns the routing key from the RoutingKey transport header.
	// RoutingKey() string
	// // RoutingDelegate returns the routing delegate from the RoutingDelegate transport header.
	// RoutingDelegate() string
	// // LocalPeer returns the local peer information for this call.
	// LocalPeer() LocalPeerInfo
	// // RemotePeer returns the remote peer information for this call.
	// RemotePeer() PeerInfo
	// // CallOptions returns a CallOptions struct suitable for forwarding a request.
	// CallOptions() *CallOptions
	// // Arg2Reader returns an ArgReader to read the second argument.
	// // The ReadCloser must be closed once the argument has been read.
	// Arg2Reader() (ArgReader, error)
	// // Arg3Reader returns an ArgReader to read the last argument.
	// // The ReadCloser must be closed once the argument has been read.
	// Arg3Reader() (ArgReader, error)
	// // Response provides access to the InboundCallResponse object which can be used
	// // to write back to the calling peer
	// Response() InboundCallResponse
}

// An InboundCallResponse is used to send the response back to the calling peer
type InboundCallResponse interface {
	// // SendSystemError returns a system error response to the peer.  The call is considered
	// // complete after this method is called, and no further data can be written.
	// SendSystemError(err error) error
	// // SetApplicationError marks the response as being an application error.  This method can
	// // only be called before any arguments have been sent to the calling peer.
	// SetApplicationError() error
	// // Blackhole indicates no response will be sent, and cleans up any resources
	// // associated with this request. This allows for services to trigger a timeout in
	// // clients without holding on to any goroutines on the server.
	// Blackhole()
	// // Arg2Writer returns a WriteCloser that can be used to write the second argument.
	// // The returned writer must be closed once the write is complete.
	// Arg2Writer() (ArgWriter, error)
	// // Arg3Writer returns a WriteCloser that can be used to write the last argument.
	// // The returned writer must be closed once the write is complete.
	// Arg3Writer() (ArgWriter, error)
}

// // Logger provides an abstract interface for logging from TChannel.
// // Applications can provide their own implementation of this interface to adapt
// // TChannel logging to whatever logging library they prefer (stdlib log,
// // logrus, go-logging, etc).  The SimpleLogger adapts to the standard go log
// // package.
// type Logger interface {
// 	// // Enabled returns whether the given level is enabled.
// 	// Enabled(level tchannel.LogLevel) bool

// 	// Fatal logs a message, then exits with os.Exit(1).
// 	Fatal(msg string)

// 	// Error logs a message at error priority.
// 	Error(msg string)

// 	// Warn logs a message at warning priority.
// 	Warn(msg string)

// 	// Infof logs a message at info priority.
// 	Infof(msg string, args ...interface{})

// 	// Info logs a message at info priority.
// 	Info(msg string)

// 	// Debugf logs a message at debug priority.
// 	Debugf(msg string, args ...interface{})

// 	// Debug logs a message at debug priority.
// 	Debug(msg string)

// 	// // Fields returns the fields that this logger contains.
// 	// Fields() tchannel.LogFields

// 	// // WithFields returns a logger with the current logger's fields and fields.
// 	// WithFields(fields ...tchannel.LogField) Logger
// }

// // StatsReporter is the the interface used to report stats.
// type StatsReporter interface {
// 	IncCounter(name string, tags map[string]string, value int64)
// 	UpdateGauge(name string, tags map[string]string, value int64)
// 	RecordTimer(name string, tags map[string]string, d time.Duration)
// }

// Peer represents a single autobahn service or client with a unique host:port.
type Peer interface {
	// HostPort returns the host:port used to connect to this peer.
	// HostPort() string
	// GetConnection(ctx context.Context) (*tchannel.Connection, error)
	// Connect(ctx context.Context) (*tchannel.Connection, error)
	// BeginCall(ctx context.Context, serviceName, methodName string, callOptions *CallOptions) (OutboundCall, error)
	// NumConnections() (inbound int, outbound int)
	// NumPendingOutbound() int
	Call(ctx context.Context, serviceName, methodName string, callOptions *CallOptions, arg, resp any) error
	RawCall(ctx context.Context, serviceName, methodName string, callOptions *CallOptions, arg []byte, resp *[]byte) error
}

type ApplicationError interface {
	error
	Type() string
	Message() string
}

// PeerList maintains a list of Peers.
type PeerList interface {
	// // SetStrategy sets customized peer selection strategy.
	// SetStrategy(sc tchannel.ScoreCalculator)
	// // Add adds a peer to the list if it does not exist, or returns any existing peer.
	// Add(hostPort string) Peer
	// GetNew(prevSelected map[string]struct{}) (Peer, error)
	// // Get returns a peer from the peer list, or nil if none can be found,
	// // will avoid previously selected peers if possible.
	// Get(prevSelected map[string]struct{}) (Peer, error)
	// // Remove removes a peer from the peer list. It returns an error if the peer cannot be found.
	// // Remove does not affect connections to the peer in any way.
	// Remove(hostPort string) error
	// GetOrAdd returns a peer for the given hostPort, creating one if it doesn't yet exist.
	GetOrAdd(hostPort string) Peer
	// // Copy returns a copy of the PeerList as a map from hostPort to peer.
	// Copy() map[string]Peer
	// // Len returns the length of the PeerList.
	// Len() int
}

// ContextWithHeaders is a Context which contains request and response headers.
type ContextWithHeaders interface {
	context.Context
	// Headers returns the call request headers.
	Headers() map[string]string
	// // ResponseHeaders returns the call response headers.
	// ResponseHeaders() map[string]string
	// // SetResponseHeaders sets the given response headers on the context.
	// SetResponseHeaders(map[string]string)
	// // Child creates a child context which stores headers separately from
	// // the parent context.
	// Child() ContextWithHeaders
}
