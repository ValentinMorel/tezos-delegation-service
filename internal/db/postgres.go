package db

import (
	"context"
	"fmt"
	"log"

	"tezos-delegation-service/config"

	_ "github.com/golang-migrate/migrate/v4/database/postgres" // Import database migration driver
	_ "github.com/golang-migrate/migrate/v4/source/file"       // Import file source for migrations
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// PgQueries wraps Queries and provides a connection pool for PostgreSQL operations
type PgQueries struct {
	*Queries               // Embedded Queries type
	pool     *pgxpool.Pool // PostgreSQL connection pool
}

// NewDatabase initializes a new PgQueries instance with a connection pool
func NewDatabase(cfg *config.Config) Querier {
	// Build the Data Source Name (DSN) for PostgreSQL connection
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUsername, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	// Create a new connection pool
	pool, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	return PgQueries{
		Queries: New(pool), // Initialize Queries with the connection pool
		pool:    pool,      // Store the connection pool
	}
}

// Connect establishes a connection pool and returns it
func Connect(cfg *config.Config) (*pgxpool.Pool, error) {
	// Build the Data Source Name (DSN) for PostgreSQL connection
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUsername, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	// Create a new connection pool
	pool, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

// InsertDelegationsBatch inserts a batch of delegations into the database
func (q PgQueries) InsertDelegationsBatch(ctx context.Context, delegations []Delegation) error {
	batch := &pgx.Batch{}
	for _, delegation := range delegations {
		batch.Queue(
			`INSERT INTO delegations (id, delegator, timestamp, amount, level)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (id, delegator, timestamp, amount, level) DO NOTHING`,
			delegation.ID, delegation.Delegator, delegation.Timestamp, delegation.Amount, delegation.Level,
		)
	}

	// Send the batch to the database
	br := q.pool.SendBatch(ctx, batch)
	defer br.Close()

	// Execute each query in the batch
	for range delegations {
		if _, err := br.Exec(); err != nil {
			return err
		}
	}

	return nil
}

// GetDelegationsByYear retrieves delegations for a specific year from the database
func (q PgQueries) GetDelegationsByYear(ctx context.Context, year int) ([]*Delegation, error) {
	var rows pgx.Rows
	var err error
	if year == 0 {
		rows, err = q.pool.Query(ctx, `
		SELECT *
		FROM delegations
		ORDER BY timestamp DESC
	`)
	} else {
		rows, err = q.pool.Query(ctx, `
		SELECT id, delegator, timestamp, amount, level
		FROM delegations
		WHERE EXTRACT(YEAR FROM timestamp) = $1
		ORDER BY timestamp DESC
	`, year)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var delegations []*Delegation
	for rows.Next() {
		var d Delegation
		if err := rows.Scan(&d.ID, &d.Delegator, &d.Timestamp, &d.Amount, &d.Level); err != nil {
			return nil, err
		}
		delegations = append(delegations, &d)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return delegations, nil
}

// DeleteDelegationsBatch deletes a batch of delegations from the database
func (q PgQueries) DeleteDelegationsBatch(ctx context.Context, delegations []Delegation) error {
	batch := &pgx.Batch{}
	for _, delegation := range delegations {
		batch.Queue("DELETE FROM delegations WHERE id = $1", delegation.ID)
	}

	// Send the batch to the database
	br := q.pool.SendBatch(ctx, batch)
	defer br.Close()

	// Execute each query in the batch
	for range delegations {
		if _, err := br.Exec(); err != nil {
			return err
		}
	}

	return nil
}
