package main

import (
	"context"
	"github.com/mis-hashemi/mahdaad-test/internal/services/course_service"
	"github.com/mis-hashemi/mahdaad-test/internal/services/dashboard_service"
	"github.com/mis-hashemi/mahdaad-test/internal/services/notification_service"
	"github.com/mis-hashemi/mahdaad-test/internal/services/search_service"
	"github.com/mis-hashemi/mahdaad-test/pkg/eventbus"
)

func main() {
	ctx := context.Background()
	bus := eventbus.New()
	defer bus.Close()

	// register subscribers
	notification_service.SubscribeEmail(ctx, bus)
	dashboard_service.SubscribeDashboard(ctx, bus)
	search_service.SubscribeIndexer(ctx, bus)

	cService := course_service.New(bus)
	cService.CreateCourse(ctx, "1", "Course 1")

}
