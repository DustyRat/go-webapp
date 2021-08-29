package rbac

import (
	"net/http"

	"github.com/dustyrat/go-webapp/internal/middleware"
)

type RBAC struct {
	// TODO: RBAC implementation
}

func (rb *RBAC) Middleware(w http.ResponseWriter, r *http.Request) (middleware.User, error) {
	// TODO: RBAC implementation
	return middleware.User{SAMAccountName: "ANONYMOUS"}, nil
}
