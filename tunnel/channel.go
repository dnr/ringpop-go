package tunnel

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"sync"

	"github.com/temporalio/ringpop-go/shared"
)

const (
	// CurrentProtocolVersion is the current version of the TChannel protocol
	// supported by this stack
	CurrentProtocolVersion = 0x03
)

type (
	Transport interface {
		Call(
			ctx context.Context,
			hostPort string,
			serviceName string,
			methodName string,
			// FIXME: headers? other options?
			req []byte,
		) ([]byte, error)

		Register(
			serviceName string,
			handler TransportHandler,
		) error
	}
	TransportHandler func(ctx context.Context, hostPort string, methodName string, req []byte) ([]byte, error)

	tun struct {
		lock        sync.Mutex
		subChannels map[string]*sub
		transport   Transport
	}

	sub struct {
		lock        sync.Mutex
		t           *tun
		serviceName string
		funcs       shared.Handlers
		onError     func(context.Context, error)
		peers       map[string]*peer
	}

	peer struct {
		t        *tun
		s        *sub
		hostPort string
	}
)

var _ shared.TChannel = (*tun)(nil)
var _ shared.SubChannel = (*sub)(nil)
var _ shared.Peer = (*peer)(nil)

func NewChannel(_ string, transport Transport) (*tun, error) {
	if transport == nil {
		// FIXME: make test one?
	}
	return &tun{
		subChannels: make(map[string]*sub),
		transport:   transport,
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

func (s *sub) Register(funcs shared.Handlers, onError func(context.Context, error)) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.funcs != nil || s.onError != nil {
		panic("duplicate registration")
	}
	s.funcs = funcs
	s.onError = onError
	return s.t.transport.Register(s.serviceName, s.handle)
}

func (s *sub) handle(ctx context.Context, hostPort string, methodName string, reqBytes []byte) (resBytes []byte, retErr error) {
	defer func() {
		if retErr != nil && s.onError != nil {
			s.onError(ctx, retErr)
		}
	}()

	// ringpop doesn't actually look at the source peer at all
	// p := s.GetOrAdd(hostPort)

	f := s.funcs[methodName]
	if f == nil {
		return nil, errors.New("unknown method")
	}

	rf := reflect.ValueOf(f)
	req := reflect.New(rf.Type().In(1))

	if err := json.Unmarshal(reqBytes, req.Interface()); err != nil {
		return nil, err
	}

	ins := []reflect.Value{reflect.ValueOf(ctx), req}
	outs := reflect.ValueOf(f).Call(ins)

	if outs[1].IsValid() && !outs[1].IsNil() {
		err, ok := outs[1].Interface().(error)
		if !ok || err == nil {
			err = errors.New("got nil error of concrete type")
		}
		return nil, err
	} else if !outs[0].IsValid() || outs[0].IsNil() {
		return nil, errors.New("invalid response")
	}

	return json.Marshal(outs[0].Interface())
}

// Peers returns the peer list for this Registrar.
func (s *sub) Peers() shared.PeerList {
	return s // this is a PeerList too
}

func (s *sub) GetOrAdd(hostPort string) shared.Peer {
	s.lock.Lock()
	defer s.lock.Unlock()
	if p, ok := s.peers[hostPort]; ok {
		return p
	}
	p := &peer{
		t:        s.t,
		s:        s,
		hostPort: hostPort,
	}
	s.peers[hostPort] = p
	return p
}

func (p *peer) Call(
	ctx context.Context, serviceName, methodName string, callOptions *shared.CallOptions, arg, resp any,
) error {
	argBytes, err := json.Marshal(arg)
	if err != nil {
		return err
	}
	var resBytes []byte
	err = p.RawCall(ctx, serviceName, methodName, callOptions, argBytes, &resBytes)
	if err != nil {
		return err
	}
	return json.Unmarshal(resBytes, resp)
}

func (p *peer) RawCall(
	ctx context.Context, serviceName, methodName string, callOptions *shared.CallOptions, arg []byte, resp *[]byte,
) error {
	// FIXME: callOptions?
	res, err := p.t.transport.Call(ctx, p.hostPort, serviceName, methodName, arg)
	if err != nil {
		return err
	}
	*resp = res
	return nil
}
