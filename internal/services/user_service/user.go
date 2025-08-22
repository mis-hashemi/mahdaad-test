package user_service

import (
	"fmt"

	"github.com/mis-hashemi/mahdaad-test/pkg/queue"
)

type Service struct {
	q *queue.SimpleQueue
}

func New(q *queue.SimpleQueue) *Service {
	return &Service{q: q}
}

func (s *Service) UpdateProfile(userID string, newData string) {
	fmt.Printf("User %s updated profile: %s\n", userID, newData)
	s.q.Publish(queue.Message{Key: userID, Data: newData})
}
