package course_service

import "github.com/mis-hashemi/mahdaad-test/pkg/eventbus"

type Service struct {
	bus *eventbus.Bus
}

func New(bus *eventbus.Bus) *Service {
	return &Service{bus: bus}
}
