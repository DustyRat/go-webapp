package file

import (
	"net/http"

	"github.com/dustyrat/go-webapp/internal/service"

	router "github.com/dustyrat/go-metrics/router/mux"

	"github.com/rs/zerolog/log"
)

func AddHandlers(r *router.Router) {
	r.Handle("/", navigation("./web/views/index.html"))
	r.Handle("/home", navigation("./web/views/home/index.html"))
	r.Handle("/user", navigation("./web/views/user/index.html"))
	r.Handle("/role", navigation("./web/views/role/index.html"))
	r.Handle("/permission", navigation("./web/views/permission/index.html"))

	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets", http.FileServer(http.Dir("./web/assets"))))
	r.PathPrefix("/template/").Handler(http.StripPrefix("/template", http.FileServer(http.Dir("./web/template"))))
	r.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger", http.FileServer(http.Dir("./swagger"))))
}

func navigation(file string) http.HandlerFunc {
	templates, err := service.Build([]string{file, "./web/template/navigation.html"})
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	return func(w http.ResponseWriter, r *http.Request) {
		service.Render(w, templates)
	}
}
