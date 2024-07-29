package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"tezos-delegation-service/internal/api"
	"tezos-delegation-service/internal/db"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockQuerier struct {
	mock.Mock
}

func (m *MockQuerier) CreateDelegation(ctx context.Context, arg db.CreateDelegationParams) (*db.CreateDelegationRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(*db.CreateDelegationRow), args.Error(1)
}

func (m *MockQuerier) GetDelegationsByYear(ctx context.Context, year int) ([]*db.Delegation, error) {
	args := m.Called(ctx, year)
	return args.Get(0).([]*db.Delegation), args.Error(1)
}

func (m *MockQuerier) InsertDelegationsBatch(ctx context.Context, delegations []db.Delegation) error {
	args := m.Called(ctx, delegations)
	return args.Error(0)
}
func (m *MockQuerier) DeleteDelegationsBatch(ctx context.Context, delegations []db.Delegation) error {
	args := m.Called(ctx, delegations)
	return args.Error(0)
}

func TestHandler_GetDelegations_Success(t *testing.T) {
	mockQuerier := new(MockQuerier)
	handler := NewHandler(mockQuerier)

	year := 2024
	expectedDelegations := []*db.Delegation{
		{ID: 1, Delegator: "Alice", Timestamp: time.Now(), Amount: 100, Level: 1},
		{ID: 2, Delegator: "Bob", Timestamp: time.Now(), Amount: 200, Level: 2},
	}

	mockQuerier.On("GetDelegationsByYear", mock.Anything, year).Return(expectedDelegations, nil)

	req, _ := http.NewRequest("GET", "/delegations?year="+strconv.Itoa(year), nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.GetXtzDelegations(c, api.GetXtzDelegationsParams{})

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"code":200`)
	assert.Contains(t, w.Body.String(), `"data":[`)
	mockQuerier.AssertExpectations(t)
}
