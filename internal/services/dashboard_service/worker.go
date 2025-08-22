package dashboard_service

import (
	"context"
	"fmt"
	"github.com/mis-hashemi/mahdaad-test/internal/entity"
	"github.com/mis-hashemi/mahdaad-test/pkg/eventbus"
)

func SubscribeDashboard(ctx context.Context, bus *eventbus.Bus) {
	bus.Subscribe("course_created", func(e eventbus.Event) {
		c := e.Data.(*entity.Course)
		fmt.Printf("Update admin dashboard for course: %s\n", c.Title)
	})
}
