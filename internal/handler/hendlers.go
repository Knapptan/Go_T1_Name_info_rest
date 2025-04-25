package handler

import (
	"effective-mobile/internal/models"
	"encoding/json"
	"net/http"
)

// "log"
// "github.com/go-chi/chi"

func CreatePerson(w http.ResponseWriter, r *http.Request) {
	var input models.PersonInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		// Валидация и возврат 400
	}
	// Логика обогощения данных 
	// Сохранение в бд
}

// 	r.Post("/", h.CreatePerson)
// 	r.Get("/", h.GetPeople)
// 	r.Put("/{id}", h.UpdatePerson)
// 	r.Delete("/{id}", h.DeletePerson)
