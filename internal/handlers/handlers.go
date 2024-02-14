package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/yablus/les30/internal/models"
	"github.com/yablus/les30/internal/requests"
)

type UserStorage interface {
	List() []*models.User
	Get(int) *models.User
	Update(int, models.User) *models.User
	Create(models.User)
	Delete(int) *models.User
}

type UserHandler struct {
	Storage UserStorage
}

func (u *UserHandler) Run() {
	r := setupServer()
	http.ListenAndServe(":8080", r)
}

func setupServer() chi.Router {
	r := chi.NewRouter()
	//models.NewStorage()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) { // GET "/"
		w.Write([]byte("OK"))
	})
	r.Mount("/users", UserRoutes())
	return r
}

func UserRoutes() chi.Router {
	r := chi.NewRouter()
	u := UserHandler{}
	r.Get("/", u.ListUsers)                // GET /users
	r.Post("/", u.CreateUser)              // POST /users
	r.Put("/{id}", u.UpdateUser)           // PUT /users/{id}
	r.Delete("/", u.DeleteUser)            // DELETE /users
	r.Post("/make_friends", u.MakeFriends) // POST /users/make_friends
	r.Get("/{id}/friends", u.GetFriends)   // GET /users/{id}/friends
	return r
}

func (u *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	//err := json.NewEncoder(w).Encode(models.ListUsers())
	err := json.NewEncoder(w).Encode(u.Storage.List())
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		log.Println("Internal error")
		return
	}
	log.Printf("List all users.")
}

func (u *UserHandler) GetFriends(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		log.Println("Internal error")
		return
	}
	//user := models.GetUser(intId)
	user := u.Storage.Get(intId)
	if user == nil {
		http.Error(w, "Not found", http.StatusNotFound)
		log.Println("Not found: Пользователь не найден")
		return
	}
	var list string
	for _, u := range models.NewStorage().Users {
		for _, v := range user.Friends {
			if u.ID == v {
				if list != "" {
					list += ", "
				}
				list += u.Name
			}
		}
	}
	wr := fmt.Sprintf("Друзья %s: %v %s", user.Name, user.Friends, list)
	log.Println("List of friends.", wr)
	w.Write([]byte(fmt.Sprint(user.Friends)))
}

func (u *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	var req requests.Create
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println("Bad request:", err.Error())
		return
	}
	//defer r.Body.Close()
	var user models.User
	models.IDs++
	user.ID = models.IDs
	user.Name = req.Name
	user.Age = req.Age
	user.Friends = req.Friends
	//models.StoreUser(user)
	u.Storage.Create(user)
	log.Printf("User created. ID=%d", user.ID)
	w.Write([]byte(fmt.Sprint(user.ID)))
}

func (u *UserHandler) MakeFriends(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated) // Здесь вернуть 201, а не 200, как указано в задании.
	var req requests.MakeFriends
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println("Bad request:", err.Error())
		return
	}
	if req.Source_id == req.Target_id {
		http.Error(w, "Bad request", http.StatusBadRequest)
		log.Println("Bad request: неверный id пользователя")
		return
	}
	if req.Source_id == 0 || req.Target_id == 0 {
		http.Error(w, "Bad request", http.StatusBadRequest)
		log.Println("Bad request: неверный id пользователя")
		return
	}
	var userS, userT models.User
	countUsers := 0
	for _, u := range models.NewStorage().Users {
		if u.ID == req.Source_id {
			userS = *u
			countUsers++
		}
		if u.ID == req.Target_id {
			userT = *u
			countUsers++
		}
	}
	if countUsers != 2 {
		http.Error(w, "Not found", http.StatusNotFound)
		log.Println("Not found: Пользователь не найден")
		return
	}
	for _, v := range userS.Friends {
		if v == userT.ID {
			http.Error(w, "Bad request", http.StatusBadRequest)
			log.Println("Bad request: Пользователи уже являются друзьями")
			return
		}
	}
	userS.Friends = append(userS.Friends, userT.ID)
	userT.Friends = append(userT.Friends, userS.ID)
	//if models.UpdateUser(userS.ID, userS) == nil || models.UpdateUser(userT.ID, userT) == nil {
	if u.Storage.Update(userS.ID, userS) == nil || u.Storage.Update(userT.ID, userT) == nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		log.Println("Internal error")
		return
	}
	wr := fmt.Sprintf("%s и %s теперь друзья", userS.Name, userT.Name)
	log.Println("Friends Added.", wr)
	w.Write([]byte(wr))
}

func (u *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var req requests.Update
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println("Bad request:", err.Error())
		return
	}
	id := chi.URLParam(r, "id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		log.Println("Internal error")
		return
	}
	//user := models.GetUser(intId)
	user := u.Storage.Get(intId)
	if user == nil {
		http.Error(w, "Not found", http.StatusNotFound)
		log.Println("Not found: Пользователь не найден")
		return
	}
	user.Age = req.NewAge
	//updatedUser := models.UpdateUser(intId, *user)
	updatedUser := u.Storage.Update(intId, *user)
	if updatedUser == nil {
		http.Error(w, "Not found", http.StatusNotFound)
		log.Println("Not found: Пользователь не найден")
		return
	}
	wr := fmt.Sprintf("Возраст %s изменен на %d", user.Name, user.Age)
	log.Println("User Updated.", wr)
	w.Write([]byte("Возраст пользователя успешно обновлен"))
}

func (u *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	var req requests.Delete
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println("Bad request:", err.Error())
		return
	}
	var user models.User
	for _, u := range models.NewStorage().Users {
		if u.ID == req.Target_id {
			user = *u
			break
		}
	}
	//if models.DeleteUser(req.Target_id) == nil {
	if u.Storage.Delete(req.Target_id) == nil {
		http.Error(w, "Not found", http.StatusNotFound)
		log.Println("Not found: Пользователь не найден")
		return
	}
	log.Printf("User deleted. Name=%s", user.Name)
	w.Write([]byte(fmt.Sprint(user.Name)))
}
