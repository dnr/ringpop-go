package tunnel

import (
	"context"
	"sync"

	"github.com/temporalio/ringpop-go/shared"
)

const (
	// CurrentProtocolVersion is the current version of the TChannel protocol
	// supported by this stack
	CurrentProtocolVersion = 0x03
)

type (
	Tunnel interface {
	}

	tun struct {
		lock        sync.Mutex
		subChannels map[string]*sub
	}

	sub struct {
		lock        sync.Mutex
		t           *tun
		serviceName string
	}

	peer struct {
	}
)

var _ shared.TChannel = (*tun)(nil)
var _ shared.SubChannel = (*sub)(nil)

func NewChannel(_ string, opts any) (*tun, error) {
	return &tun{
		subChannels: make(map[string]*sub),
	}, nil
}

func (t *tun) GetSubChannel(serviceName string) shared.SubChannel {
	t.lock.Lock()
	defer t.lock.Unlock()
	if sub, ok := t.subChannels[serviceName]; ok {
		return sub
	}
	sub := &sub{
		t:           t,
		serviceName: serviceName,
	}
	t.subChannels[serviceName] = sub
	return sub
}

func (t *tun) PeerInfo() shared.LocalPeerInfo {
	panic("unimpl")
}

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
	return s.serviceName
}

// Register registers a handler for ServiceName and the given method.
func (s *sub) Register(h shared.Handler, methodName string) {
	panic("unimpl")
}

// Peers returns the peer list for this Registrar.
func (s *sub) Peers() shared.PeerList {
	return s // this is a PeerList too
}

func (s *sub) GetOrAdd(hostPort string) shared.Peer {
	panic("unimpl")
}

func (p *peer) Call(
	ctx context.Context, serviceName, methodName string, callOptions *shared.CallOptions, arg, resp any,
) error {
	panic("unimpl")
}

// Handlers is the map from method names to handlers.
type Handlers map[string]interface{}

// JsonRegister registers the specified methods specified as a map from method name to the
// JSON handler function. The handler functions should have the following signature:
// func(context.Context, *ArgType)(*ResType, error)
func JsonRegister(registrar shared.Registrar, funcs Handlers, onError func(context.Context, error)) error {
	// FIXME: what about onError?
	for name, handler := range funcs {
		registrar.Register(handler, name)
	}
	return nil
}
