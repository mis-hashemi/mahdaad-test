package profile_sync

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/mis-hashemi/mahdaad-test/pkg/queue"
)

func Handler(msg queue.Message) error {
	// simulate flaky external service
	if rand.Intn(3) == 0 {
		return fmt.Errorf("profile service temporarily unavailable")
	}
	fmt.Printf("Synced profile for %s with data: %+v\n", msg.Key, msg.Data)
	return nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
