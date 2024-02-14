package run

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/yablus/les30/internal/handlers"
)

func Run() {
	r := setupServer()
	http.ListenAndServe(":3000", r)
}

func setupServer() chi.Router {
	r := chi.NewRouter()
	//models.NewStorage()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) { // GET "/"
		w.Write([]byte("OK"))
	})
	r.Mount("/users", UserRoutes())
	http.ListenAndServe(":8080", r)
	return r
}

func UserRoutes() chi.Router {
	r := chi.NewRouter()
	u := handlers.UserHandler{}
	r.Get("/", u.ListUsers)                // GET /users
	r.Post("/", u.CreateUser)              // POST /users
	r.Put("/{id}", u.UpdateUser)           // PUT /users/{id}
	r.Delete("/", u.DeleteUser)            // DELETE /users
	r.Post("/make_friends", u.MakeFriends) // POST /users/make_friends
	r.Get("/{id}/friends", u.GetFriends)   // GET /users/{id}/friends
	return r
}
