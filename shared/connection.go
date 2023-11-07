package shared

// --- copied from tchannel-go/connection.go

import "fmt"

// PeerVersion contains version related information for a specific peer.
// These values are extracted from the init headers.
type PeerVersion struct {
	Language        string `json:"language"`
	LanguageVersion string `json:"languageVersion"`
	TChannelVersion string `json:"tchannelVersion"`
}

// PeerInfo contains information about a TChannel peer
type PeerInfo struct {
	// The host and port that can be used to contact the peer, as encoded by net.JoinHostPort
	HostPort string `json:"hostPort"`

	// The logical process name for the peer, used for only for logging / debugging
	ProcessName string `json:"processName"`

	// IsEphemeral returns whether the remote host:port is ephemeral (e.g. not listening).
	IsEphemeral bool `json:"isEphemeral"`

	// Version returns the version information for the remote peer.
	Version PeerVersion `json:"version"`
}

func (p PeerInfo) String() string {
	return fmt.Sprintf("%s(%s)", p.HostPort, p.ProcessName)
}

// IsEphemeralHostPort returns whether the connection is from an ephemeral host:port.
func (p PeerInfo) IsEphemeralHostPort() bool {
	return p.IsEphemeral
}

// LocalPeerInfo adds service name to the peer info, only required for the local peer.
type LocalPeerInfo struct {
	PeerInfo

	// ServiceName is the service name for the local peer.
	ServiceName string `json:"serviceName"`
}

func (p LocalPeerInfo) String() string {
	return fmt.Sprintf("%v: %v", p.ServiceName, p.PeerInfo)
}
