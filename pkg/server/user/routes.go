package user

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"net/http"
	"testZaShtat/pkg/errors"
)

type HttpUserHandler struct {
	repository Repository
}

func Routes(repo Repository) *chi.Mux {
	handler := &HttpUserHandler{repo}
	router := chi.NewRouter()
	router.Group(func(router chi.Router) {
		router.Post("/", handler.createUser())
		router.Post("/login", handler.login())
	})
	router.Group(func(router chi.Router) {
		router.Use(jwtauth.Verifier(TokenAuth))
		router.Use(jwtauth.Authenticator)
		router.Get("/", handler.getAllUsers())
		router.Get("/{userID}", handler.getUser())
		router.Put("/{userID}", handler.modifyUser())
		router.Delete("/{userID}", handler.deleteUser())
	})
	return router
}

func (h *HttpUserHandler) login() http.HandlerFunc {

	type loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	type loginResponse struct {
		Token string `json:"token"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &loginRequest{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			render.Render(w, r, errors.ErrInvalidRequest(err))
			return
		}
		user, err := h.repository.GetByUsername(req.Username)
		if err != nil {
			render.Render(w, r, errors.ErrInvalidRequest(fmt.Errorf("username or password is invalid")))
			return
		}
		err = user.CheckPassword(req.Password)
		if err != nil {
			render.Render(w, r, errors.ErrInvalidRequest(fmt.Errorf("username or password is invalid")))
			return
		}
		token, err := user.GenToken()
		if err != nil {
			render.Render(w, r, errors.ErrInvalidRequest(err))
			return
		}
		response := &loginResponse{
			Token: token,
		}
		render.Status(r, http.StatusCreated)
		render.JSON(w, r, response)
	}
}

func (h *HttpUserHandler) getUser() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userID")
		if userID == "" {
			render.Render(w, r, errors.ErrInvalidRequest(fmt.Errorf("invalid userID")))
			return
		}
		user, err := h.repository.GetByID(userID)
		if err != nil {
			render.Render(w, r, errors.ErrInvalidRequest(err))
			return
		}
		render.JSON(w, r, user)
		render.Status(r, http.StatusOK)
	}
}

func (h *HttpUserHandler) getAllUsers() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		users, err := h.repository.GetAllUsers()
		if err != nil {
			render.Render(w, r, errors.ErrInvalidRequest(err))
			return
		}
		render.Status(r, http.StatusOK)
		render.JSON(w, r, users)
	}
}

func (h *HttpUserHandler) deleteUser() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userID")
		if userID == "" {
			render.Render(w, r, errors.ErrInvalidRequest(fmt.Errorf("invalid userID")))
			return
		}
		err := h.repository.Delete(userID)
		if err != nil {
			render.Render(w, r, errors.ErrInvalidRequest(err))
			return
		}
		render.NoContent(w, r)
	}
}

func (h *HttpUserHandler) createUser() http.HandlerFunc {

	type createRequest struct {
		Username   string `json:"username"`
		FirstName  string `json:"firstName"`
		MiddleName string `json:"middleName"`
		LastName   string `json:"lastName"`
		Password   string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &createRequest{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			render.Render(w, r, errors.ErrInvalidRequest(err))
			return
		}
		user := &User{
			Username:   req.Username,
			FirstName:  req.FirstName,
			MiddleName: req.MiddleName,
			LastName:   req.LastName,
		}
		if err := user.SetPassword(req.Password); err != nil {
			render.Render(w, r, errors.ErrInvalidRequest(err))
			return
		}
		if err := h.repository.Store(user); err != nil {
			render.Render(w, r, errors.ErrInvalidRequest(err))
			return
		}
		render.Status(r, http.StatusCreated)
		render.JSON(w, r, user)
	}
}

func (h *HttpUserHandler) modifyUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := &User{}
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			render.Render(w, r, errors.ErrInvalidRequest(err))
			return
		}
		err = h.repository.Update(user)
		if err != nil {
			render.Render(w, r, errors.ErrInvalidRequest(err))
			return
		}
		render.Status(r, http.StatusCreated)
		render.JSON(w, r, user)
	}
}
