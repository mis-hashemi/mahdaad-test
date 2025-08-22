package saga

import (
	"context"
	"github.com/mis-hashemi/mahdaad-test/pkg/retry"
)

// SagaStep represents one unit of work in the saga
type SagaStep struct {
	Name         string
	Action       func(ctx context.Context) error
	Compensation func(ctx context.Context) error
	Retry        *retry.Retry // optional retry policy
}
