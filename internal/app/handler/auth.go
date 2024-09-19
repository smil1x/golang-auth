package handler

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"golang-auth/internal/app/dto"
	"net/http"
)

func (h *Handlers) GetTokens(w http.ResponseWriter, r *http.Request) {

	guid := r.URL.Query().Get("guid")
	if guid == "" {
		h.error(w, http.StatusBadRequest, errors.New("guid is missing"))
		return
	}

	userIp := h.service.Authorization.GetIp(r)
	tokens, err := h.service.Authorization.GetTokens(guid, userIp)

	if err != nil {
		h.error(w, http.StatusNotFound, err)
		return
	}
	h.respond(w, http.StatusOK, tokens)
}

func (h *Handlers) TokensRefreshing(w http.ResponseWriter, r *http.Request) {
	tokens := &dto.TokensDTO{}
	if err := json.NewDecoder(r.Body).Decode(tokens); err != nil {
		h.error(w, http.StatusBadRequest, err)
		return
	}

	validate := validator.New()
	if err := validate.Struct(tokens); err != nil {
		h.error(w, http.StatusBadRequest, err)
		return
	}

	userIp := h.service.Authorization.GetIp(r)
	tokens, err := h.service.RefreshTokens(tokens, userIp)
	if err != nil {
		h.error(w, http.StatusBadRequest, err)
		return
	}
	h.respond(w, http.StatusOK, tokens)
}

func (h *Handlers) error(w http.ResponseWriter, code int, err error) {
	h.respond(w, code, map[string]string{"error": err.Error()})
}

func (h *Handlers) respond(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		}
	}
}
