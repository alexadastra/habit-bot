package inmemory

import (
	"context"
	"log"
	"sync"

	"github.com/alexadastra/habit_bot/internal/models"
)

const logPrefix = "[InMemoryStateStorage]"

type InMemoryStateStorage struct {
	statesByIDs *sync.Map
}

func NewInMemoryStateStorage() *InMemoryStateStorage {
	return &InMemoryStateStorage{
		statesByIDs: &sync.Map{},
	}
}

func (s *InMemoryStateStorage) FetchByID(ctx context.Context, userID int64) (models.UserState, error) {
	state, ok := s.statesByIDs.Load(userID)
	if !ok {
		log.Printf("%s state not found for userID %d! return default state", logPrefix, userID)
		return models.DefaultUserState, nil
	}

	st, ok := state.(models.UserState)
	if !ok {
		log.Printf("%s cannot convert found state '%v' to string for userID %d! return default state", logPrefix, state, userID)
		return models.DefaultUserState, nil
	}

	return st, nil
}
func (s *InMemoryStateStorage) SetByID(ctx context.Context, userID int64, newState models.UserState) error {
	s.statesByIDs.Store(userID, newState)
	return nil
}
