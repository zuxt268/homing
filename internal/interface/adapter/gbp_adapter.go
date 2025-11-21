package adapter

import "context"

type GbpAdapter interface {
	GetCustomer(ctx context.Context, accountID, locationID string) string
}

type gbpAdapter struct {
}

func NewGbpAdapter() GbpAdapter {
	return &gbpAdapter{}
}

func (a *gbpAdapter) GetCustomer(ctx context.Context, accountID, locationID string) string {
	return ""
}
