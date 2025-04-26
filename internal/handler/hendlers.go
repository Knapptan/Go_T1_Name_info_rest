package handler

import (
	"effective-mobile/internal/models"
	"encoding/json"
	"net/http"
)

// "log"
// "github.com/go-chi/chi"

// r.Post("/", h.CreatePerson)
func CreatePerson(w http.ResponseWriter, r *http.Request) {
	var input models.PersonInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		// Валидация и возврат 400
	}
	// Логика обогощения данных
	// Сохранение в бд
}

// r.Get("/", h.GetPeople)
func GetPeople() {

}

// r.Put("/{id}", h.UpdatePerson)
func UpdatePerson() {

}

// r.Delete("/{id}", h.DeletePerson)
func DeletePerson() {

}
