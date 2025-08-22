package course_service

import (
	"context"
	"github.com/mis-hashemi/mahdaad-test/internal/entity"
	"github.com/mis-hashemi/mahdaad-test/pkg/eventbus"
)

func (s *Service) CreateCourse(ctx context.Context, id, title string) *entity.Course {
	c := &entity.Course{ID: id, Title: title}
	s.bus.Publish(eventbus.Event{
		Name: "course_created",
		Data: c,
	})
	return c
}
