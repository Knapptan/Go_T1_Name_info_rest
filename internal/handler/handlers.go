package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"effective-mobile/internal/models"
	"effective-mobile/internal/service"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

type Handler struct {
	svc    *service.PersonService
	logger *zap.Logger
}

func NewHandler(svc *service.PersonService, logger *zap.Logger) *Handler {
	return &Handler{
		svc:    svc,
		logger: logger.With(zap.String("component", "http_handler")),
	}
}

// Post Handler Создание обогащённого человека
func (h *Handler) CreatePerson(w http.ResponseWriter, r *http.Request) {

	start := time.Now()
	logMethod := zap.String("method", r.Method)
	logPath := zap.String("path", r.URL.Path)

	defer func() {
		h.logger.Debug("Request completed",
			logMethod,
			logPath,
			zap.Duration("duration", time.Since(start)))
	}()

	var input models.PersonInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.logger.Error("Invalid request body", logMethod, logPath, zap.Error(err))
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if input.Name == "" || input.Surname == "" {
		h.logger.Warn("Validation failed", logMethod, logPath, zap.String("reason", "missing required fields"),
			zap.Any("input", input))
		http.Error(w, "name and surname are required", http.StatusBadRequest)
		return
	}

	person, err := h.svc.CreatePersonEnriched(r.Context(), input)
	if err != nil {
		h.logger.Error("Failed to create person", logMethod, logPath, zap.Error(err), zap.Any("input", input))
		http.Error(w, "failed to create person: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(person); err != nil {
		h.logger.Error("Failed to encode response", logMethod, logPath, zap.Error(err))
		return
	}

	h.logger.Info("Person created successfully",
		logMethod, logPath,
		zap.Int("person_id", person.ID),
		zap.String("name", person.Name),
	)
}

// Get Handler Получение списка людей с фильтрами и пагинацией
func (h *Handler) GetPeople(w http.ResponseWriter, r *http.Request) {
	filters := models.Filters{
		Name:   r.URL.Query().Get("name"),
		Gender: r.URL.Query().Get("gender"),
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}

	filters.Offset = (page - 1) * limit
	filters.Limit = limit

	people, err := h.svc.GetPersons(r.Context(), filters)
	if err != nil {
		http.Error(w, "failed to fetch people: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(people)
}

// Put Handler Обновление данных человека
func (h *Handler) UpdatePerson(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var upd models.PersonUpdate
	if err := json.NewDecoder(r.Body).Decode(&upd); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := h.svc.UpdatePerson(r.Context(), id, upd); err != nil {
		http.Error(w, "failed to update person: "+err.Error(), http.StatusInternalServerError)
		return
	}

}

// Delete Handler Удаление человека по ID
func (h *Handler) DeletePerson(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	err = h.svc.DeletePerson(r.Context(), id)
	if err != nil {

		if strings.Contains(err.Error(), "not found") {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, "failed to delete person: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
