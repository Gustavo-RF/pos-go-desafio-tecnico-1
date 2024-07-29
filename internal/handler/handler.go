package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Gustavo-RF/desafio-tecnico-1/internal/entities"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	response := entities.Response{
		Message: "success",
	}
	json.NewEncoder(w).Encode(response)
}
