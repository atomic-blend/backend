package rspamdservice

import (
	rspamdclient "github.com/atomic-blend/backend/shared/services/rspamd/client"
	rspamdinterfaces "github.com/atomic-blend/backend/shared/services/rspamd/interfaces"
)

// Wrapper wraps the existing rspamd functionality
type Wrapper struct {
	client *rspamdclient.Client
}

// NewRspamdService creates a new rspamd service wrapper
func NewRspamdService() rspamdinterfaces.RspamdServiceInterface {
	config := rspamdclient.DefaultConfig()
	client := rspamdclient.NewClient(config)
	return &Wrapper{
		client: client,
	}
}

// NewRspamdServiceWithConfig creates a new rspamd service wrapper with custom config
func NewRspamdServiceWithConfig(config *rspamdclient.Config) rspamdinterfaces.RspamdServiceInterface {
	client := rspamdclient.NewClient(config)
	return &Wrapper{
		client: client,
	}
}

// CheckMessage sends a message to Rspamd for spam checking
func (r *Wrapper) CheckMessage(req *rspamdclient.CheckRequest) (*rspamdclient.CheckResponse, error) {
	return r.client.CheckMessage(req)
}

// Ping sends a ping request to check if Rspamd is available
func (r *Wrapper) Ping() error {
	return r.client.Ping()
}

// NewClient creates a new Rspamd client with the given configuration
func (r *Wrapper) NewClient(config *rspamdclient.Config) *rspamdclient.Client {
	return rspamdclient.NewClient(config)
}

// DefaultConfig returns default configuration with environment variable support
func (r *Wrapper) DefaultConfig() *rspamdclient.Config {
	return rspamdclient.DefaultConfig()
}
