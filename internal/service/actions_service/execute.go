package actions_service

import (
	"context"
	"log"

	"github.com/alexadastra/habit_bot/internal/models"
)

func (s *ActionsService) execute(ctx context.Context, action models.Action) error {
	log.Printf("doing action '%s' with the title '%s'!", action.ID, action.Name)

	if action, ok := s.actions[action.Name]; ok {
		return action(ctx)
	}

	log.Printf("action not found by name '%s'", action.Name)

	return nil
}
