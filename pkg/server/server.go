package server

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	_ "github.com/lib/pq"
	config "github.com/spf13/viper"
	"log"
	"net"
	"net/http"
	"testZaShtat/pkg/server/user"
	"testZaShtat/pkg/store"
)

func initConfig() {
	config.SetConfigFile("config.yaml")
	err := config.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func Routes(userRepo *user.UserRepository) *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.Logger,
		middleware.DefaultCompress,
		middleware.RedirectSlashes,
		middleware.Recoverer,
	)

	router.Route("/v1", func(r chi.Router) {
		r.Mount("/api/users", user.Routes(userRepo))
	})

	return router
}

func ListenAndServe() {
	initConfig()
	db, err := store.NewDB()
	if err != nil {
		log.Panic(err)
	}
	userRepository := &user.UserRepository{DB: db}
	router := Routes(userRepository)
	host := net.JoinHostPort(config.GetString("server.host"), config.GetString("server.port"))
	log.Fatal(http.ListenAndServe(host, router))
}
