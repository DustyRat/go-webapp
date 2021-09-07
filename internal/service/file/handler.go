package file

import (
	"net/http"

	"github.com/dustyrat/go-webapp/internal/service"

	router "github.com/dustyrat/go-metrics/router/mux"

	"github.com/rs/zerolog/log"
)

func AddHandlers(r *router.Router) {
	r.Handle("/", index("./web/static/index.html"))
	r.Handle("/home", home("./web/static/home/index.html"))
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets", http.FileServer(http.Dir("./web/assets"))))
	r.PathPrefix("/template/").Handler(http.StripPrefix("/template", http.FileServer(http.Dir("./web/template"))))
	r.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger", http.FileServer(http.Dir("./swagger"))))
}

func index(file string) http.HandlerFunc {
	templates, err := service.Build([]string{file, "./web/template/navigation.html"})
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	return func(w http.ResponseWriter, r *http.Request) {
		service.Render(w, templates)
	}
}

func home(file string) http.HandlerFunc {
	templates, err := service.Build([]string{file, "./web/template/navigation.html"})
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	return func(w http.ResponseWriter, r *http.Request) {
		service.Render(w, templates)
	}
}
