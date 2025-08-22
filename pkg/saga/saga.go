package saga

import (
	"context"
	"fmt"
	"github.com/mis-hashemi/mahdaad-test/pkg/logger"
)

// Saga orchestrates multi-step workflows with rollback
type Saga struct {
	steps  []SagaStep
	logger *logger.Logger
}

func New(l *logger.Logger) *Saga {
	return &Saga{
		logger: l,
	}
}

// AddStep appends a saga step
func (s *Saga) AddStep(step SagaStep) {
	s.steps = append(s.steps, step)
}

// Execute runs the saga flow
func (s *Saga) Execute(ctx context.Context) error {
	var completed []SagaStep

	for _, step := range s.steps {
		s.logger.Info("executing step", "step", step.Name)

		var err error
		if step.Retry != nil {
			err = step.Retry.Do(ctx, func() error { return step.Action(ctx) })
		} else {
			err = step.Action(ctx)
		}

		if err != nil {
			s.logger.Error("step failed", "step", step.Name, "err", err)

			// Rollback in reverse order
			for i := len(completed) - 1; i >= 0; i-- {
				compStep := completed[i]
				if compStep.Compensation != nil {
					s.logger.Info("rolling back step", "step", compStep.Name)
					if cerr := compStep.Compensation(ctx); cerr != nil {
						s.logger.Error("compensation failed",
							"step", compStep.Name, "err", cerr)
					}
				}
			}
			return fmt.Errorf("saga failed at step %s: %w", step.Name, err)
		}

		completed = append(completed, step)
	}

	s.logger.Info("saga completed successfully")
	return nil
}
