package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/spf13/cobra"

	"gpb.ru/hr/internal/hr/repos/postgres"
	"gpb.ru/hr/internal/hr/services"
)

func Server() *cobra.Command {
	pgurl := ""

	cmd := &cobra.Command{
		Use:   "serve [address]",
		Short: "Run HR API server on the given address.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			repos, err := postgres.New(pgurl)
			if err != nil {
				log.Printf("[error] database connection error: %s", err)
				return
			}
			server := services.NewServer(args[0], repos.Candidate, repos.Vacancy)

			done := make(chan struct{})
			go func() {
				defer close(done)
				err := server.Run()
				if err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Printf("[error] server running error: %s", err)
				}
			}()

			select {
			case <-cmd.Context().Done():
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				err := server.Close(ctx)
				if err != nil {
					log.Printf("[error] server shutdown error: %s", err)
				}
			case <-done:
			}
		},
	}

	cmd.Flags().StringVar(&pgurl, "db", "localhost:5432", "Postgres url.")

	return cmd
}
