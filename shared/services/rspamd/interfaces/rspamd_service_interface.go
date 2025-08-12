// Package rspamdinterfaces contains the interfaces for the Rspamd service
package rspamdinterfaces

import rspamdclient "github.com/atomic-blend/backend/shared/services/rspamd/client"

// RspamdServiceInterface defines the interface for rspamd operations
type RspamdServiceInterface interface {
	CheckMessage(req *rspamdclient.CheckRequest) (*rspamdclient.CheckResponse, error)
	Ping() error
	NewClient(config *rspamdclient.Config) *rspamdclient.Client
	DefaultConfig() *rspamdclient.Config
}
