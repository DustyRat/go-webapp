package middleware

import (
	"net/http"

	"github.com/DustyRat/go-webapp/internal/service"

	"github.com/gorilla/context"
)

const key = "User"

type User struct {
	Cn             string   `json:"cn"`
	Company        string   `json:"company"`
	GivenName      string   `json:"givenName"`
	Initials       string   `json:"initials"`
	Mail           string   `json:"mail"`
	MemberOf       []string `json:"memberOf"`
	SAMAccountName string   `json:"sAMAccountName"`
	Sn             string   `json:"sn"`
	Title          string   `json:"title"`
}

func RBAC(fn func(w http.ResponseWriter, r *http.Request) (User, error), next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writer := &writer{ResponseWriter: w}
		user, err := fn(writer, r)
		if err != nil {
			service.RespondWithError(w, http.StatusBadGateway, err)
			return
		}

		// If the writer has already been written to return
		if writer.statusCode > 0 {
			return
		}

		context.Set(r, key, user)
		next(w, r)
	})
}

func GetUser(r *http.Request) User {
	if user, ok := context.Get(r, key).(User); ok {
		return user
	}
	return User{}
}
