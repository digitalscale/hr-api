package repos

import (
	"context"

	"github.com/google/uuid"

	"gpb.ru/hr/internal/hr/entities"
)

type CandidateRepo interface {
	GetByID(context.Context, uuid.UUID) (*entities.Candidate, error)
	List(context.Context, uuid.UUID) ([]entities.Candidate, error)
	Create(context.Context, *entities.Candidate) error
	Update(context.Context, *entities.Candidate) error
}
