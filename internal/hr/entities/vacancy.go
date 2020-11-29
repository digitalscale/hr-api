package entities

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type VacancyStatus byte

const (
	VacancyStatusNone VacancyStatus = iota
	VacancyStatusDraft
	VacancyStatusActive
	VacancyStatusInactive
	vacancyStatusCount
)

var vacancyStatusStrings = []string{
	"none",
	"draft",
	"active",
	"inactive",
}

func (status VacancyStatus) String() string {
	if status >= vacancyStatusCount {
		return vacancyStatusStrings[VacancyStatusNone]
	}
	return vacancyStatusStrings[status]
}

func (status VacancyStatus) MarshalText() ([]byte, error) {
	v := status.String()
	return []byte(v), nil
}

func (status *VacancyStatus) Scan(src interface{}) error {
	switch v := src.(type) {
	case string:
		return status.UnmarshalText([]byte(v))
	case []byte:
		return status.UnmarshalText(v)
	}
	return nil
}

var vacancyStatusTexts = map[string]VacancyStatus{
	"":         VacancyStatusNone,
	"none":     VacancyStatusNone,
	"draft":    VacancyStatusDraft,
	"active":   VacancyStatusActive,
	"inactive": VacancyStatusInactive,
}

var ErrInvalidVacancyStatus = errors.New("invalid vacancy status")

func (status *VacancyStatus) UnmarshalText(data []byte) error {
	v, ok := vacancyStatusTexts[string(data)]
	if !ok {
		return ErrInvalidVacancyStatus
	}
	*status = v
	return nil
}

type Skill struct {
	Title     string `json:"title"`
	Important bool   `json:"important"`
}

type Vacancy struct {
	ID           uuid.UUID     `json:"id"`
	TemplateID   uuid.UUID     `json:"templateID"`
	Title        string        `json:"title"`
	Status       VacancyStatus `json:"status"`
	Area         string        `json:"area,omitempty"`
	Department   string        `json:"department,omitempty"`
	Skills       []Skill       `json:"skills"`
	Duties       []string      `json:"duties"`
	Requirements []string      `json:"requirements"`
	Experience   uint32        `json:"experience"`
	Created      time.Time     `json:"created"`
	Updated      time.Time     `json:"updated"`
}

func (v *Vacancy) Validate() error {
	return nil
}
