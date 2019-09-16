package user

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
)

type HttpUserHandler struct {
	repository Repository
}

func Routes(repo Repository) *chi.Mux {
	handler := &HttpUserHandler{repo}
	router := chi.NewRouter()
	router.Post("/", handler.createUser)
	router.Get("/", handler.GetAllUsers)
	router.Get("/{userID}", handler.GetUser)
	router.Put("/{userID}", handler.modifyUser)
	router.Delete("/{userID}", handler.deleteUser)
	return router
}

func (h *HttpUserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	id, err := strconv.ParseInt(userID, 10, 32)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	user, err := h.repository.GetByID(rune(id))
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	data, err := json.Marshal(user)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(data)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *HttpUserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.repository.GetAllUsers()
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	data, err := json.Marshal(users)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *HttpUserHandler) deleteUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	id, err := strconv.ParseInt(userID, 10, 32)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	err = h.repository.Delete(rune(id))
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *HttpUserHandler) createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	user := &User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	err = h.repository.Store(user)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *HttpUserHandler) modifyUser(w http.ResponseWriter, r *http.Request) {
	user := &User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	err = h.repository.Update(user)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
