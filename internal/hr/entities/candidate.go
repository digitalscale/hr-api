package entities

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Gender byte

const (
	GenderNone Gender = iota
	GenderMale
	GenderFemale
	genderCount
)

var genderStrings = []string{
	"none",
	"male",
	"female",
}

func (gender Gender) String() string {
	if gender >= genderCount {
		return genderStrings[GenderNone]
	}
	return genderStrings[gender]
}

func (gender Gender) MarshalText() ([]byte, error) {
	v := gender.String()
	return []byte(v), nil
}

var genderTexts = map[string]Gender{
	"":       GenderNone,
	"none":   GenderNone,
	"male":   GenderMale,
	"female": GenderFemale,
}

var ErrInvalidGender = errors.New("invalid gender")

func (gender *Gender) UnmarshalText(data []byte) error {
	v, ok := genderTexts[string(data)]
	if !ok {
		return ErrInvalidGender
	}
	*gender = v
	return nil
}

func (gender *Gender) Scan(src interface{}) error {
	switch v := src.(type) {
	case string:
		return gender.UnmarshalText([]byte(v))
	case []byte:
		return gender.UnmarshalText(v)
	}
	return nil
}

type EducationLevel byte

const (
	EducationLevelNone EducationLevel = iota
	EducationLevelSecondary
	EducationLevelSpecialSecondary
	EducationLevelUnfinishedHigher
	EducationLevelHigher
	EducationLevelBachelor
	EducationLevelMaster
	EducationLevelCandidate
	EducationLevelDoctor
	educationLevelCount
)

var educationLevelStrings = []string{
	"none",
	"secondary",
	"specialSecondary",
	"unfinishedHigher",
	"higher",
	"bachelor",
	"master",
	"candidate",
	"doctor",
}

func (lvl EducationLevel) String() string {
	if lvl >= educationLevelCount {
		return educationLevelStrings[EducationLevelNone]
	}
	return educationLevelStrings[lvl]
}

func (lvl EducationLevel) MarshalText() ([]byte, error) {
	v := lvl.String()
	return []byte(v), nil
}

var educationLevelTexts = map[string]EducationLevel{
	"":                 EducationLevelNone,
	"none":             EducationLevelNone,
	"secondary":        EducationLevelSecondary,
	"specialSecondary": EducationLevelSpecialSecondary,
	"unfinishedHigher": EducationLevelUnfinishedHigher,
	"higher":           EducationLevelHigher,
	"bachelor":         EducationLevelBachelor,
	"master":           EducationLevelMaster,
	"candidate":        EducationLevelCandidate,
	"doctor":           EducationLevelDoctor,
}

var ErrInvalidEducationLevel = errors.New("invalid education level")

func (lvl *EducationLevel) UnmarshalText(data []byte) error {
	v, ok := educationLevelTexts[string(data)]
	if !ok {
		return ErrInvalidEducationLevel
	}
	*lvl = v
	return nil
}

func (lvl *EducationLevel) Scan(src interface{}) error {
	switch v := src.(type) {
	case string:
		return lvl.UnmarshalText([]byte(v))
	case []byte:
		return lvl.UnmarshalText(v)
	}
	return nil
}

type Education struct {
	Title string `json:"title"`
	Year  uint32 `json:"year"`
}

type Experience struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Start       *time.Time `json:"start"`
	End         *time.Time `json:"end"`
}

type Candidate struct {
	ID             uuid.UUID      `json:"id"`
	Name           string         `json:"name"`
	Phone          string         `json:"phone"`
	Email          string         `json:"email"`
	Specialization string         `json:"specialization"`
	Gender         Gender         `json:"gender"`
	BirthDate      *time.Time     `json:"birthDate"`
	Area           string         `json:"area"`
	Salary         uint32         `json:"salary"`
	EducationLevel EducationLevel `json:"educationLevel"`
	Education      []Education    `json:"education"`
	Experience     []Experience   `json:"experience"`
	Languages      []string       `json:"languages"`
	Skills         []string       `json:"skills"`
	Created        time.Time      `json:"created"`
	Updated        time.Time      `json:"updated"`
}
