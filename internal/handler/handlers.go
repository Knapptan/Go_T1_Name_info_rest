// @Title Effective Mobile API
// @Version 1.0
// @Description REST API для работы с данными людей

// @ContactName API Support
// @ContactEmail support@effective-mobile.ru

// @Host localhost:8080
// @BasePath /api/v1

// @SecurityDefinitions.basic BasicAuth

// Handlers.go
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

// CreatePerson godoc
// @Summary Создать нового человека
// @Description Создание новой персоны с обогащением данных
// @Tags people
// @Accept  json
// @Produce json
// @Param   input body models.PersonInput true "Данные человека"
// @Success 201 {object} models.PersonEnriched
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /people [post]
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

// GetPeople godoc
// @Summary Получить список людей
// @Description Получение списка людей с фильтрацией и пагинацией
// @Tags people
// @Produce json
// @Param   name    query string false "Фильтр по имени"
// @Param   gender  query string false "Фильтр по полу"
// @Param   page    query int    false "Номер страницы" default(1)
// @Param   limit   query int    false "Лимит на странице" default(10)
// @Success 200 {array} models.PersonEnriched
// @Failure 500 {object} map[string]string
// @Router /people [get]
func (h *Handler) GetPeople(w http.ResponseWriter, r *http.Request) {
	logMethod := zap.String("method", r.Method)
	logPath := zap.String("path", r.URL.Path)

	start := time.Now()
	defer func() {
		h.logger.Debug("Request completed",
			logMethod,
			logPath,
			zap.Duration("duration", time.Since(start)))
	}()

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

	h.logger.Debug("Fetching people",
		logMethod,
		logPath,
		zap.Any("filters", filters),
		zap.Int("page", page),
		zap.Int("limit", limit),
	)

	people, err := h.svc.GetPersons(r.Context(), filters)
	if err != nil {
		h.logger.Error("Failed to fetch people",
			logMethod,
			logPath,
			zap.Error(err),
			zap.Any("filters", filters),
		)
		http.Error(w, "failed to fetch people: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(people); err != nil {
		h.logger.Error("Failed to encode response",
			logMethod,
			logPath,
			zap.Error(err),
		)
		return
	}

	h.logger.Info("People fetched successfully",
		logMethod,
		logPath,
		zap.Int("count", len(people)),
	)
}

// UpdatePerson godoc
// @Summary Обновить данные человека
// @Description Обновление данных существующей персоны
// @Tags people
// @Accept json
// @Produce json
// @Param id path int true "ID человека"
// @Param update body models.PersonUpdate true "Данные для обновления"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /people/{id} [patch]
func (h *Handler) UpdatePerson(w http.ResponseWriter, r *http.Request) {
	logMethod := zap.String("method", r.Method)
	logPath := zap.String("path", r.URL.Path)

	start := time.Now()
	defer func() {
		h.logger.Debug("Request completed",
			logMethod,
			logPath,
			zap.Duration("duration", time.Since(start)))
	}()

	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		h.logger.Warn("Invalid ID parameter",
			logMethod,
			logPath,
			zap.String("id_param", idParam),
			zap.Error(err),
		)
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var upd models.PersonUpdate
	if err := json.NewDecoder(r.Body).Decode(&upd); err != nil {
		h.logger.Error("Invalid request body",
			logMethod,
			logPath,
			zap.Error(err),
		)
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	h.logger.Debug("Updating person",
		logMethod,
		logPath,
		zap.Int("person_id", id),
		zap.Any("update_data", upd),
	)

	if err := h.svc.UpdatePerson(r.Context(), id, upd); err != nil {
		h.logger.Error("Failed to update person",
			logMethod,
			logPath,
			zap.Error(err),
			zap.Int("person_id", id),
			zap.Any("update_data", upd),
		)
		http.Error(w, "failed to update person: "+err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info("Person updated successfully",
		logMethod,
		logPath,
		zap.Int("person_id", id),
	)
	w.WriteHeader(http.StatusOK)
}

// DeletePerson godoc
// @Summary Удалить человека
// @Description Удаление персоны по идентификатору
// @Tags people
// @Produce json
// @Param id path int true "ID человека"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /people/{id} [delete]
func (h *Handler) DeletePerson(w http.ResponseWriter, r *http.Request) {
	logMethod := zap.String("method", r.Method)
	logPath := zap.String("path", r.URL.Path)

	start := time.Now()
	defer func() {
		h.logger.Debug("Request completed",
			logMethod,
			logPath,
			zap.Duration("duration", time.Since(start)))
	}()

	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		h.logger.Warn("Invalid ID parameter",
			logMethod,
			logPath,
			zap.String("id_param", idParam),
			zap.Error(err),
		)
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	h.logger.Debug("Deleting person",
		logMethod,
		logPath,
		zap.Int("person_id", id),
	)

	err = h.svc.DeletePerson(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			h.logger.Warn("Person not found",
				logMethod,
				logPath,
				zap.Int("person_id", id),
			)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		h.logger.Error("Failed to delete person",
			logMethod,
			logPath,
			zap.Error(err),
			zap.Int("person_id", id),
		)
		http.Error(w, "failed to delete person: "+err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info("Person deleted successfully",
		logMethod,
		logPath,
		zap.Int("person_id", id),
	)
	w.WriteHeader(http.StatusNoContent)
}
