package file

import (
	"log"
	"net/http"
	"text/template"

	router "github.com/dustyrat/go-metrics/router/mux"
)

func AddHandlers(r *router.Router) {
	r.Handle("/", index())
	r.Handle("/home", home())
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets", http.FileServer(http.Dir("./web/assets"))))
	r.PathPrefix("/template/").Handler(http.StripPrefix("/template", http.FileServer(http.Dir("./web/template"))))
	r.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger", http.FileServer(http.Dir("./swagger"))))
}

func index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Initialize a slice containing the paths to the two files. Note that the
		// home.page.tmpl file must be the *first* file in the slice.
		files := []string{
			"./web/static/index.html",
			"./web/template/navigation.html",
		}

		// Use the template.ParseFiles() function to read the files and store the
		// templates in a template set. Notice that we can pass the slice of file paths
		// as a variadic parameter?
		ts, err := template.ParseFiles(files...)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error", 500)
			return
		}

		err = ts.Execute(w, nil)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error", 500)
		}
	}
}

func home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Initialize a slice containing the paths to the two files. Note that the
		// home.page.tmpl file must be the *first* file in the slice.
		files := []string{
			"./web/static/home/index.html",
			"./web/template/navigation.html",
		}

		// Use the template.ParseFiles() function to read the files and store the
		// templates in a template set. Notice that we can pass the slice of file paths
		// as a variadic parameter?
		ts, err := template.ParseFiles(files...)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error", 500)
			return
		}

		err = ts.Execute(w, nil)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error", 500)
		}
	}
}
