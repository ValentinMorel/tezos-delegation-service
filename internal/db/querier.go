package db

import (
	"context"
)

type Querier interface {
	CreateDelegation(ctx context.Context, arg CreateDelegationParams) (*CreateDelegationRow, error)
	GetDelegationsByYear(ctx context.Context, year int) ([]*Delegation, error)
	InsertDelegationsBatch(ctx context.Context, delegations []Delegation) error
	DeleteDelegationsBatch(ctx context.Context, delegations []Delegation) error
}
