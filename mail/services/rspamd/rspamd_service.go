package rspamdservice

import (
	rspamdclient "github.com/atomic-blend/backend/mail/services/rspamd/client"
	rspamdinterfaces "github.com/atomic-blend/backend/mail/services/rspamd/interfaces"
)

// RspamdServiceWrapper wraps the existing rspamd functionality
type RspamdServiceWrapper struct {
	client *rspamdclient.Client
}

// NewRspamdService creates a new rspamd service wrapper
func NewRspamdService() rspamdinterfaces.RspamdServiceInterface {
	config := rspamdclient.DefaultConfig()
	client := rspamdclient.NewClient(config)
	return &RspamdServiceWrapper{
		client: client,
	}
}

// NewRspamdServiceWithConfig creates a new rspamd service wrapper with custom config
func NewRspamdServiceWithConfig(config *rspamdclient.Config) rspamdinterfaces.RspamdServiceInterface {
	client := rspamdclient.NewClient(config)
	return &RspamdServiceWrapper{
		client: client,
	}
}

// CheckMessage sends a message to Rspamd for spam checking
func (r *RspamdServiceWrapper) CheckMessage(req *rspamdclient.CheckRequest) (*rspamdclient.CheckResponse, error) {
	return r.client.CheckMessage(req)
}

// Ping sends a ping request to check if Rspamd is available
func (r *RspamdServiceWrapper) Ping() error {
	return r.client.Ping()
}

// NewClient creates a new Rspamd client with the given configuration
func (r *RspamdServiceWrapper) NewClient(config *rspamdclient.Config) *rspamdclient.Client {
	return rspamdclient.NewClient(config)
}

// DefaultConfig returns default configuration with environment variable support
func (r *RspamdServiceWrapper) DefaultConfig() *rspamdclient.Config {
	return rspamdclient.DefaultConfig()
}
