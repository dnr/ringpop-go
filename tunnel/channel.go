package tunnel

import (
	"context"

	"github.com/temporalio/ringpop-go/shared"
)

const (
	// CurrentProtocolVersion is the current version of the TChannel protocol
	// supported by this stack
	CurrentProtocolVersion = 0x03
)

type (
	tun struct {
	}

	sub struct {
	}
)

var _ shared.TChannel = (*tun)(nil)
var _ shared.SubChannel = (*sub)(nil)

func NewChannel(serviceName string, opts any) (*tun, error) {
	panic("unimpl")
}

func (t *tun) GetSubChannel(fixme string) shared.SubChannel {
	panic("unimpl")
}

func (t *tun) PeerInfo() shared.LocalPeerInfo {
	panic("unimpl")
}

// test only
func (t *tun) ListenAndServe(hostport string) error {
	panic("unimpl")
}

// test only
func (t *tun) Ping(context.Context, string) error {
	panic("unimpl")
}

func (t *tun) Close() {
	panic("unimpl")
}

func (s *sub) ServiceName() string {
	panic("unimpl")
}

// Register registers a handler for ServiceName and the given method.
func (s *sub) Register(h shared.Handler, methodName string) {
	panic("unimpl")
}

// // Logger returns the logger for this Registrar.
// func (s *sub) Logger() shared.Logger {
// 	panic("unimpl")
// }

// // StatsReporter returns the stats reporter for this Registrar
// func (s *sub) StatsReporter() shared.StatsReporter {
// 	panic("unimpl")
// }

// // StatsTags returns the tags that should be used.
// func (s *sub) StatsTags() map[string]string {
// 	panic("unimpl")
// }

// Peers returns the peer list for this Registrar.
func (s *sub) Peers() shared.PeerList {
	panic("unimpl")
}

// JsonCallPeer makes a JSON call using the given peer.
func JsonCallPeer(ctx shared.ContextWithHeaders, peer shared.Peer, serviceName, method string, arg, resp interface{}) error {
	panic("unimpl")
}

// Handlers is the map from method names to handlers.
type Handlers map[string]interface{}

// JsonRegister registers the specified methods specified as a map from method name to the
// JSON handler function. The handler functions should have the following signature:
// func(context.Context, *ArgType)(*ResType, error)
func JsonRegister(registrar shared.Registrar, funcs Handlers, onError func(context.Context, error)) error {
	panic("unimpl")
}

// RawWriteArgs writes the given arguments to the call, and returns the response args.
func RawWriteArgs(call shared.OutboundCall, arg2, arg3 []byte) ([]byte, []byte, shared.OutboundCallResponse, error) {
	panic("unimpl")
}
