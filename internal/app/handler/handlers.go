package handler

import (
	"github.com/gorilla/mux"
	"golang-auth/internal/app/service"
	"net/http"
)

type Handlers struct {
	service *service.Service
}

func NewHandlers(service *service.Service) *Handlers {
	return &Handlers{
		service: service,
	}
}

func (h *Handlers) InitHandlers() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/auth/tokens", h.GetTokens).Methods(http.MethodGet)
	r.HandleFunc("/auth/refreshing", h.TokensRefreshing).Methods(http.MethodPost)

	return r
}
