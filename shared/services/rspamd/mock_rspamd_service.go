// Package rspamdservice contains the mock Rspamd service
package rspamdservice

import (
	rspamdclient "github.com/atomic-blend/backend/shared/services/rspamd/client"
	rspamdinterfaces "github.com/atomic-blend/backend/shared/services/rspamd/interfaces"
	"github.com/stretchr/testify/mock"
)

// MockRspamdService provides a mock implementation of rspamd service
type MockRspamdService struct {
	mock.Mock
}

// Ensure MockRspamdService implements the interface
var _ rspamdinterfaces.RspamdServiceInterface = (*MockRspamdService)(nil)

// CheckMessage sends a message to Rspamd for spam checking
func (m *MockRspamdService) CheckMessage(req *rspamdclient.CheckRequest) (*rspamdclient.CheckResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*rspamdclient.CheckResponse), args.Error(1)
}

// Ping sends a ping request to check if Rspamd is available
func (m *MockRspamdService) Ping() error {
	args := m.Called()
	return args.Error(0)
}

// NewClient creates a new Rspamd client with the given configuration
func (r *MockRspamdService) NewClient(config *rspamdclient.Config) *rspamdclient.Client {
	args := r.Called(config)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*rspamdclient.Client)
}

// DefaultConfig returns default configuration with environment variable support
func (r *MockRspamdService) DefaultConfig() *rspamdclient.Config {
	args := r.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*rspamdclient.Config)
}
