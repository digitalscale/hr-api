package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"

	"gpb.ru/hr/internal/hr/entities"
)

type VacancyRepo struct {
	db *pgxpool.Pool
}

func NewVacancyRepo(pool *pgxpool.Pool) *VacancyRepo {
	return &VacancyRepo{db: pool}
}

var ErrVacancyNotFound = errors.New("vacancy not found")

func (repo *VacancyRepo) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*entities.Vacancy, error) {
	vacancyRows, err := repo.db.Query(
		ctx,
		`SELECT * FROM vacancy.vacancy WHERE id = $1`,
		id.String(),
	)
	if err != nil {
		return nil, err
	}
	defer vacancyRows.Close()

	if !vacancyRows.Next() {
		return nil, ErrVacancyNotFound
	}

	var vacancy entities.Vacancy
	err = vacancyRows.Scan(
		&vacancy.ID,
		&vacancy.TemplateID,
		&vacancy.Title,
		&vacancy.Status,
		&vacancy.Area,
		&vacancy.Department,
		&vacancy.Duties,
		&vacancy.Requirements,
		&vacancy.Experience,
		&vacancy.Created,
		&vacancy.Updated,
	)
	if err != nil {
		return nil, err
	}

	skillRows, err := repo.db.Query(
		ctx,
		`SELECT * FROM vacancy.skill WHERE vacancy_id = $1`,
		id,
	)
	if err != nil {
		return nil, err
	}
	defer skillRows.Close()

	for skillRows.Next() {
		skill := entities.Skill{}
		err = skillRows.Scan(
			nil,
			&skill.Title,
			&skill.Important,
		)
		if err != nil {
			return nil, err
		}
		vacancy.Skills = append(vacancy.Skills, skill)
	}

	return &vacancy, nil
}

func (repo *VacancyRepo) Create(
	ctx context.Context,
	vacancy *entities.Vacancy,
) error {
	tx, err := repo.db.Begin(ctx)
	if err != nil {
		return err
	}

	vacancy.ID = uuid.New()
	vacancy.Created = time.Now()
	vacancy.Updated = time.Now()

	_, err = tx.Exec(
		ctx,
		`INSERT INTO vacancy.vacancy VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		vacancy.ID,
		vacancy.TemplateID,
		vacancy.Title,
		vacancy.Status.String(),
		vacancy.Area,
		vacancy.Department,
		vacancy.Duties,
		vacancy.Requirements,
		vacancy.Experience,
		vacancy.Created,
		vacancy.Updated,
	)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	for _, skill := range vacancy.Skills {
		_, err = tx.Exec(
			ctx,
			`INSERT INTO vacancy.skill VALUES($1,$2,$3)`,
			&vacancy.ID,
			&skill.Title,
			&skill.Important,
		)
		if err != nil {
			tx.Rollback(ctx)
			return err
		}
	}

	return tx.Commit(ctx)
}

func (repo *VacancyRepo) Update(
	ctx context.Context,
	vacancy *entities.Vacancy,
) error {

	tx, err := repo.db.Begin(ctx)
	if err != nil {
		return err
	}

	vacancy.Updated = time.Now()

	_, err = tx.Exec(
		ctx,
		`
			UPDATES vacancy.vacancy SET
				template_id = $2,
				title = $3,
				status = $4,
				area = $5,
				department = $6,
				duties = $7,
				requirements = $8,
				experience = $9,
				updated = $10
			WHERE id = $1
		`,
		vacancy.ID,
		vacancy.TemplateID,
		vacancy.Title,
		vacancy.Status.String(),
		vacancy.Area,
		vacancy.Department,
		vacancy.Duties,
		vacancy.Requirements,
		vacancy.Experience,
		vacancy.Updated,
	)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	tx.Exec(ctx, `DELETE FROM vacancy.skill WHERE vacancy_id = $1`, vacancy.ID)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	for _, skill := range vacancy.Skills {
		_, err = tx.Exec(
			ctx,
			`INSERT INTO vacancy.skill VALUES($1,$2,$3)`,
			&vacancy.ID,
			&skill.Title,
			&skill.Important,
		)
		if err != nil {
			tx.Rollback(ctx)
			return err
		}
	}

	return tx.Commit(ctx)
}

func (repo *VacancyRepo) List(ctx context.Context) ([]entities.Vacancy, error) {
	vacancyRows, err := repo.db.Query(
		ctx,
		`SELECT * FROM vacancy.vacancy ORDER BY updated desc`,
	)
	if err != nil {
		return nil, err
	}
	defer vacancyRows.Close()

	vacancies := make([]entities.Vacancy, 0, 1000)
	index := make(map[uuid.UUID]int)
	for vacancyRows.Next() {
		vacancy := entities.Vacancy{}
		err = vacancyRows.Scan(
			&vacancy.ID,
			&vacancy.TemplateID,
			&vacancy.Title,
			&vacancy.Status,
			&vacancy.Area,
			&vacancy.Department,
			&vacancy.Duties,
			&vacancy.Requirements,
			&vacancy.Experience,
			&vacancy.Created,
			&vacancy.Updated,
		)
		if err != nil {
			return nil, err
		}
		vacancies = append(vacancies, vacancy)
		index[vacancy.ID] = len(vacancies) - 1
	}

	if len(vacancies) > 0 {
		args := make([]interface{}, 0, len(vacancies))
		params := make([]string, 0, len(vacancies))
		i := 1
		for id := range index {
			args = append(args, id)
			params = append(params, fmt.Sprintf("$%d", i))
			i++
		}
		query := fmt.Sprintf(
			"SELECT * FROM vacancy.skill WHERE vacancy_id in (%s)",
			strings.Join(params, ","),
		)

		skillRows, err := repo.db.Query(ctx, query, args...)
		if err != nil {
			return nil, err
		}
		defer skillRows.Close()

		for skillRows.Next() {
			var vacancyID uuid.UUID
			skill := entities.Skill{}
			err = skillRows.Scan(
				&vacancyID,
				&skill.Title,
				&skill.Important,
			)
			if err != nil {
				return nil, err
			}
			id := index[vacancyID]
			vacancies[id].Skills = append(vacancies[id].Skills, skill)
		}
	}

	return vacancies, nil
}
