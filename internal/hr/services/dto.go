package services

import (
	"time"

	"github.com/google/uuid"

	"gpb.ru/hr/internal/hr/entities"
)

type Vacancy struct {
	ID         uuid.UUID              `json:"id"`
	Title      string                 `json:"title"`
	Status     entities.VacancyStatus `json:"status"`
	Area       string                 `json:"area"`
	Department string                 `json:"department"`
	Created    time.Time              `json:"created"`
	Updated    time.Time              `json:"updated"`
}

type ListVacanciesResponse struct {
	Items []Vacancy `json:"items"`
	Token string    `json:"token,omitempty"`
}

type ErrorResponse struct {
	Code int    `json:"code"`
	Text string `json:"message"`
}
