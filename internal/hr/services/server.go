package services

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"gpb.ru/hr/internal/hr/entities"
	"gpb.ru/hr/internal/hr/repos"
)

// Server defines how the HR API interacta and stores its state.
type Server struct {
	mu     sync.Mutex
	server *http.Server

	candidate repos.CandidateRepo
	vacancy   repos.VacancyRepo
}

// NewServer creates new server with the given properties.
func NewServer(
	addr string,
	candidate repos.CandidateRepo,
	vacancy repos.VacancyRepo,
) *Server {

	server := &Server{
		candidate: candidate,
		vacancy:   vacancy,
	}

	router := mux.NewRouter()
	router.HandleFunc("/vacancies", server.ListVacancies).Methods(http.MethodGet)
	router.HandleFunc("/vacancies/{id}", server.GetVacancy).Methods(http.MethodGet)
	router.HandleFunc("/vacancies", server.CreateVacancy).Methods(http.MethodPost)
	router.HandleFunc("/vacancies/{id}", server.UpdateVacancy).Methods(http.MethodPost)

	router.HandleFunc("/cards", server.ListCards).Methods(http.MethodGet)
	router.HandleFunc("/cards/{id}", server.GetCard).Methods(http.MethodGet)
	router.HandleFunc("/cards/{id}", server.MoveCard).Methods(http.MethodPut)

	server.server = &http.Server{
		Addr:    addr,
		Handler: WithCORS(router),
	}
	return server
}

// ListVacancies return a list of vacancies.
func (srv *Server) ListVacancies(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	result, err := srv.vacancy.List(req.Context())
	if err != nil {
		log.Printf("[error] [server] error listing vacancies: %s", err)
		return
	}

	items := make([]Vacancy, len(result))
	for i, vacancy := range result {
		items[i] = Vacancy{
			ID:         vacancy.ID,
			Title:      vacancy.Title,
			Status:     vacancy.Status,
			Area:       vacancy.Area,
			Department: vacancy.Department,
			Created:    vacancy.Created,
			Updated:    vacancy.Updated,
		}
	}

	token := ""
	if len(items) > 0 {
		token = items[len(items)-1].ID.String()
	}

	response := ListVacanciesResponse{
		Items: items,
		Token: token,
	}
	err = writeJSON(w, http.StatusOK, response)
	if err != nil {
		log.Printf("[error] [server] error listing vacancies: %s", err)
	}
}

// GetVacancy returns detailed information about specified vacancy.
func (srv *Server) GetVacancy(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	id, ok := mux.Vars(req)["id"]
	if !ok {
		err := writeJSON(w, http.StatusNotFound, errors.New("not found"))
		if err != nil {
			log.Printf("[error] [server] error get vacancy: %s", err)
		}
		return
	}

	vacancyID, err := uuid.Parse(id)
	if err != nil {
		log.Printf("[error] [server] error get vacancy: %s", err)
		err := writeError(w, http.StatusBadRequest, err)
		if err != nil {
			log.Printf("[error] [server] error get vacancy: %s", err)
		}
		return
	}

	response, err := srv.vacancy.GetByID(req.Context(), vacancyID)
	if err != nil {
		log.Printf("[error] [server] error get vacancy: %s", err)
		err := writeError(w, http.StatusBadRequest, err)
		if err != nil {
			log.Printf("[error] [server] error get vacancy: %s", err)
		}
		return
	}

	err = writeJSON(w, http.StatusOK, response)
	if err != nil {
		log.Printf("[error] [server] %s", err)
	}
}

// CreateVacancy creates vacancy with the given properties.
func (srv *Server) CreateVacancy(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	var vacancy entities.Vacancy
	err := json.NewDecoder(req.Body).Decode(&vacancy)
	if err != nil {
		log.Printf("[error] [server] error creating vacancy: %s", err)
		writeError(w, http.StatusBadRequest, err)
		return
	}

	err = vacancy.Validate()
	if err != nil {
		log.Printf("[error] [server] error creating vacancy: %s", err)
		writeError(w, http.StatusBadRequest, err)
		return
	}

	for i := 0; i < 5; i++ {
		err = srv.vacancy.Create(req.Context(), &vacancy)
		if err != nil {
			continue
		}
		break
	}

	if err != nil {
		log.Printf("[error] [server] error creating vacancy: %s", err)
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	err = writeJSON(w, http.StatusOK, vacancy)
	if err != nil {
		log.Printf("[error] [server] error creating vacancy: %s", err)
	}
}

// UpdateVacancy updates properties of the given vacancy.
func (srv *Server) UpdateVacancy(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	id, ok := mux.Vars(req)["id"]
	if !ok {
		err := writeJSON(w, http.StatusNotFound, errors.New("not found"))
		if err != nil {
			log.Printf("[error] [server] %s", err)
		}
		return
	}

	vacancyID, err := uuid.Parse(id)
	if err != nil {
		log.Printf("[error] [server] error updating vacancy: %s", err)
		err := writeError(w, http.StatusBadRequest, err)
		if err != nil {
			log.Printf("[error] [server] error updating vacancy: %s", err)
		}
		return
	}

	var vacancy entities.Vacancy
	err = json.NewDecoder(req.Body).Decode(&vacancy)
	if err != nil {
		log.Printf("[error] [server] error updating vacancy: %s", err)
		writeError(w, http.StatusBadRequest, err)
		return
	}
	vacancy.ID = vacancyID

	for i := 0; i < 5; i++ {
		err = srv.vacancy.Update(req.Context(), &vacancy)
		if err != nil {
			continue
		}
		break
	}

	if err != nil {
		log.Printf("[error] [server] error updating vacancy: %s", err)
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	err = writeJSON(w, http.StatusOK, vacancy)
	if err != nil {
		log.Printf("[error] [server] error updating vacancy: %s", err)
	}
}

// ListCards return a list of canban cards with the given filter.
func (srv *Server) ListCards(w http.ResponseWriter, req *http.Request) {
}

// GetCard returns detailed information about specified canban card.
func (srv *Server) GetCard(w http.ResponseWriter, req *http.Request) {
}

// MoveCard moves card to the specified column.
func (srv *Server) MoveCard(w http.ResponseWriter, req *http.Request) {
}

// AddComment adds comments to the specified card.
func (srv *Server) AddComment(w http.ResponseWriter, req *http.Request) {
}

// Run runs the server on the given address.
func (srv *Server) Run() error {
	listener, err := net.Listen("tcp", srv.server.Addr)
	if err != nil {
		log.Printf("[error] [server] %s", err)
		return err
	}
	log.Printf("[info] [server] listen on %s", srv.server.Addr)
	return srv.server.Serve(listener)
}

// Close gracefully stops the sserver.
func (srv *Server) Close(ctx context.Context) error {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	if srv.server != nil {
		return srv.server.Shutdown(ctx)
	}

	return nil
}

func writeJSON(w http.ResponseWriter, code int, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(response)
}

func writeError(w http.ResponseWriter, code int, err error) error {
	return writeJSON(w, code, ErrorResponse{Code: code, Text: err.Error()})
}
