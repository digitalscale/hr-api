package postgres

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"

	"gpb.ru/hr/internal/hr/repos"
)

type Postgres struct {
	pool      *pgxpool.Pool
	Candidate repos.CandidateRepo
	Vacancy   repos.VacancyRepo
}

func New(uri string) (*Postgres, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	config, err := pgxpool.ParseConfig(uri)
	if err != nil {
		return nil, err
	}
	config.ConnConfig.Logger = &logger{}
	log.Printf(
		"[info] [postgres] connecting to %s %s",
		config.ConnConfig.Host,
		config.ConnConfig.Database,
	)

	pool, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return &Postgres{
		pool:    pool,
		Vacancy: NewVacancyRepo(pool),
	}, nil
}

func (pg *Postgres) Close(ctx context.Context) error {
	done := make(chan struct{})
	go func() {
		defer close(done)
		pg.pool.Close()
	}()

	select {
	case <-ctx.Done():
	case <-done:
	}

	return nil
}
