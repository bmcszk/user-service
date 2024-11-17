package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/bmcszk/user-service/logic"
)

const defaultLimit = 10

type Handler struct {
	http.Handler
	service *logic.Service
}

func NewHandler(service *logic.Service) *Handler {
	router := http.NewServeMux()
	h := &Handler{
		Handler: router,
		service: service,
	}
	router.HandleFunc("POST /users", h.createUser)
	router.HandleFunc("GET /users/{id}", h.getUserByID)
	router.HandleFunc("PUT /users/{id}", h.updateUserByID)
	router.HandleFunc("DELETE /users/{id}", h.deleteUserByID)
	router.HandleFunc("GET /users", h.listUsers)
	return h
}

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	var user logic.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		handleInputError(w, err)
		return
	}
	res, err := h.service.CreateUser(r.Context(), user)
	if err != nil {
		handleLogicError(w, err)
		return
	}
	handleResult(w, http.StatusCreated, res)
}

func (h *Handler) getUserByID(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		handleInputError(w, err)
		return
	}
	user, err := h.service.GetUserByID(r.Context(), id)
	if err != nil {
		handleLogicError(w, err)
		return
	}
	handleResult(w, http.StatusOK, user)
}

func (h *Handler) updateUserByID(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		handleInputError(w, err)
		return
	}
	var user logic.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		handleInputError(w, err)
		return
	}
	res, err := h.service.UpdateUserByID(r.Context(), id, user)
	if err != nil {
		handleLogicError(w, err)
		return
	}
	handleResult(w, http.StatusOK, res)
}

func (h *Handler) deleteUserByID(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		handleInputError(w, err)
		return
	}
	if err := h.service.DeleteUserByID(r.Context(), id); err != nil {
		handleLogicError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) listUsers(w http.ResponseWriter, r *http.Request) {
	limit, err := getParam(r, "limit", defaultLimit)
	if err != nil {
		handleInputError(w, err)
		return
	}
	offset, err := getParam(r, "offset", 0)
	if err != nil {
		handleInputError(w, err)
		return
	}
	users, err := h.service.ListUsers(r.Context(), limit, offset)
	if err != nil {
		handleLogicError(w, err)
		return
	}
	handleResult(w, http.StatusOK, users)
}

func getParam(r *http.Request, key string, defaultValue int32) (int32, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return defaultValue, nil
	}
	i, err := strconv.ParseInt(param, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("parsing int %s: %w", key, err)
	}
	return int32(i), nil
}

func handleResult(w http.ResponseWriter, code int, v any) {
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func handleInputError(w http.ResponseWriter, err error) {
	code := http.StatusBadRequest
	slog.With("error", err, "code", code).Error("invalid input")
	handleResult(w, code, ApiError{
		StatusCode: code,
		Message:    err.Error(),
	})
}

func handleLogicError(w http.ResponseWriter, err error) {
	code := getStatusCode(err)
	slog.With("error", err, "code", code).Error("logic error")
	handleResult(w, code, ApiError{
		StatusCode: code,
		Message:    err.Error(),
	})
}

func getStatusCode(err error) int {
	if errors.Is(err, logic.ErrUserNotFound) {
		return http.StatusNotFound
	}
	if errors.Is(err, logic.ErrUserAlreadyExists) {
		return http.StatusConflict
	}
	if errors.Is(err, logic.ErrUserNameEmpty) {
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}
