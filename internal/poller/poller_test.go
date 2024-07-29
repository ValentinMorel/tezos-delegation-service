package poller

import (
	"context"
	"testing"
	"tezos-delegation-service/config"
	"tezos-delegation-service/internal/db"
	"tezos-delegation-service/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockQuerier struct {
	mock.Mock
}

func (m MockQuerier) CreateDelegation(ctx context.Context, arg db.CreateDelegationParams) (*db.CreateDelegationRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(*db.CreateDelegationRow), args.Error(1)
}

func (m MockQuerier) GetDelegationsByYear(ctx context.Context, year int) ([]*db.Delegation, error) {
	args := m.Called(ctx, year)
	return args.Get(0).([]*db.Delegation), args.Error(1)
}

func (m MockQuerier) InsertDelegationsBatch(ctx context.Context, delegations []db.Delegation) error {
	args := m.Called(ctx, delegations)
	return args.Error(0)
}
func (m MockQuerier) DeleteDelegationsBatch(ctx context.Context, delegations []db.Delegation) error {
	args := m.Called(ctx, delegations)
	return args.Error(0)
}

func TestNewPoller(t *testing.T) {
	cfg := &config.Config{
		PollingInterval: "10s",
		TezosAPIURL:     "https://api.tzkt.io/v1/operations/delegations",
	}
	querier := MockQuerier{} // You will need to implement this mock

	p := NewPoller(cfg, querier, logger.NewLogger())

	assert.Equal(t, cfg, p.cfg)
	assert.Equal(t, querier, p.querier)
	assert.NotNil(t, p.httpClient)
	assert.NotNil(t, p.cache)
}
