package repos

import (
	"context"

	"github.com/google/uuid"

	"gpb.ru/hr/internal/hr/entities"
)

type VacancyRepo interface {
	GetByID(context.Context, uuid.UUID) (*entities.Vacancy, error)
	List(context.Context) ([]entities.Vacancy, error)
	Create(context.Context, *entities.Vacancy) error
	Update(context.Context, *entities.Vacancy) error
}
