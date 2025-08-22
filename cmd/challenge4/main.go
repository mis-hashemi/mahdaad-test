package main

import (
	"context"
	"time"

	"github.com/mis-hashemi/mahdaad-test/internal/services/profile_sync"
	"github.com/mis-hashemi/mahdaad-test/internal/services/user_service"
	"github.com/mis-hashemi/mahdaad-test/pkg/queue"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	q := queue.New(10)

	// start consumer for syncing profile
	q.StartConsumer(ctx, profile_sync.Handler)

	// simulate user service
	uService := user_service.New(q)
	uService.UpdateProfile("u1", "new email")
	uService.UpdateProfile("u2", "new phone")

	// wait a bit to see retries/success
	time.Sleep(10 * time.Second)
}
