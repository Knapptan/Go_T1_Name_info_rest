package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"effective-mobile/internal/models"
	"effective-mobile/internal/service"
	// "github.com/go-chi/chi"
	// "github.com/go-chi/chi/v5"
)

type Handler struct {
	svc *service.PersonService
}

func NewHandler(svc *service.PersonService) *Handler {
	return &Handler{svc: svc}
}

// r.Post("/", h.CreatePerson)
func (h *Handler) CreatePerson(w http.ResponseWriter, r *http.Request) {
	var input models.PersonInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		log.Printf("invalid JSON body: %v", err)
		return
	}

	if input.Name == "" || input.Surname == "" {
		http.Error(w, "name and surname are required", http.StatusBadRequest)
		return
	}

	person, err := h.svc.CreatePersonEnriched(r.Context(), input)
	if err != nil {
		http.Error(w, "failed to create person: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(person)

}

// r.Get("/", h.GetPeople)
func GetPeople(w http.ResponseWriter, r *http.Request) {

}

// r.Put("/{id}", h.UpdatePerson)
func UpdatePerson(w http.ResponseWriter, r *http.Request) {

}

// r.Delete("/{id}", h.DeletePerson)
func DeletePerson(w http.ResponseWriter, r *http.Request) {

}
