package search_service

import (
	"context"
	"fmt"
	"github.com/mis-hashemi/mahdaad-test/internal/entity"
	"github.com/mis-hashemi/mahdaad-test/pkg/eventbus"
)

func SubscribeIndexer(ctx context.Context, bus *eventbus.Bus) {
	bus.Subscribe("course_created", func(e eventbus.Event) {
		c := e.Data.(*entity.Course)
		fmt.Printf("Index course in search system: %s\n", c.Title)
	})
}
