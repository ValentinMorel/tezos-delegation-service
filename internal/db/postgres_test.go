package db

import (
	"context"
	"testing"
	"tezos-delegation-service/config"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDatabase(t *testing.T) {
	cfg := &config.Config{
		DBUsername: "username",
		DBPassword: "password",
		DBHost:     "localhost",
		DBPort:     "5432",
		DBName:     "db_test",
	}

	db := NewDatabase(cfg)
	pgQueries, ok := db.(PgQueries)
	require.True(t, ok)
	require.NotNil(t, pgQueries.pool)
	assert.IsType(t, &pgxpool.Pool{}, pgQueries.pool)
}

func TestConnect(t *testing.T) {
	cfg := &config.Config{
		DBUsername: "username",
		DBPassword: "password",
		DBHost:     "localhost",
		DBPort:     "5432",
		DBName:     "db_test",
	}

	pool, err := Connect(cfg)
	require.NoError(t, err)
	require.NotNil(t, pool)
	assert.IsType(t, &pgxpool.Pool{}, pool)
}

func TestInsertDelegationsBatch(t *testing.T) {
	// Setup
	cfg := &config.Config{
		DBUsername: "username",
		DBPassword: "password",
		DBHost:     "localhost",
		DBPort:     "5432",
		DBName:     "db_test",
	}

	// Connect
	db := NewDatabase(cfg)
	pgQueries, ok := db.(PgQueries)
	require.True(t, ok)
	ctx := context.Background()

	// Prepare
	delegations := []Delegation{
		{ID: 1, Delegator: "Alex", Timestamp: time.Now(), Amount: 100, Level: 1},
		{ID: 2, Delegator: "Bob", Timestamp: time.Now(), Amount: 200, Level: 2},
	}

	// Test
	err := pgQueries.InsertDelegationsBatch(ctx, delegations)
	require.NoError(t, err)

	// Verify
	insertedDelegations, err := pgQueries.GetDelegationsByYear(ctx, time.Now().Year())
	require.NoError(t, err)
	assert.Len(t, insertedDelegations, 2)

	err = pgQueries.DeleteDelegationsBatch(ctx, delegations)
	require.NoError(t, err)

}

func TestGetDelegationsByYear(t *testing.T) {
	// Setup
	cfg := &config.Config{
		DBUsername: "username",
		DBPassword: "password",
		DBHost:     "localhost",
		DBPort:     "5432",
		DBName:     "db_test",
	}

	// Connect
	db := NewDatabase(cfg)
	pgQueries, ok := db.(PgQueries)
	require.True(t, ok)
	ctx := context.Background()

	// Insert test data
	currentYear := time.Now().Year()
	delegations := []Delegation{
		{ID: 1, Delegator: "Alex", Timestamp: time.Now(), Amount: 100, Level: 1},
		{ID: 2, Delegator: "Bob", Timestamp: time.Now(), Amount: 200, Level: 2},
	}
	err := pgQueries.InsertDelegationsBatch(ctx, delegations)
	require.NoError(t, err)

	// Test retrieval
	result, err := pgQueries.GetDelegationsByYear(ctx, time.Now().Year())
	require.NoError(t, err)
	assert.Len(t, result, 2)
	for _, d := range result {
		assert.Equal(t, currentYear, d.Timestamp.Year())
	}

	err = pgQueries.DeleteDelegationsBatch(ctx, delegations)
	require.NoError(t, err)
}
