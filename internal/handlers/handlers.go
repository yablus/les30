package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
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
	Storage models.Storage
}

func (uh *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(uh.Storage.List())
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		log.Println("Internal error")
		return
	}
	log.Printf("List all users.")
}

func (uh *UserHandler) GetFriends(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		log.Println("Internal error")
		return
	}
	user := uh.Storage.Get(intId)
	if user == nil {
		http.Error(w, "Not found", http.StatusNotFound)
		log.Println("Not found: Пользователь не найден")
		return
	}
	var list string
	for _, u := range uh.Storage.Users {
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

func (uh *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	var req requests.Create
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println("Bad request:", err.Error())
		return
	}
	var user models.User
	models.IDs++
	user.ID = models.IDs
	user.Name = req.Name
	user.Age = req.Age
	user.Friends = req.Friends
	uh.Storage.Create(user)
	log.Printf("User created. ID=%d", user.ID)
	w.Write([]byte(fmt.Sprint(user.ID)))
}

func (uh *UserHandler) MakeFriends(w http.ResponseWriter, r *http.Request) {
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
	for _, u := range uh.Storage.Users {
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
	if uh.Storage.Update(userS.ID, userS) == nil || uh.Storage.Update(userT.ID, userT) == nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		log.Println("Internal error")
		return
	}
	wr := fmt.Sprintf("%s и %s теперь друзья", userS.Name, userT.Name)
	log.Println("Friends Added.", wr)
	w.Write([]byte(wr))
}

func (uh *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
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
	user := uh.Storage.Get(intId)
	if user == nil {
		http.Error(w, "Not found", http.StatusNotFound)
		log.Println("Not found: Пользователь не найден")
		return
	}
	user.Age = req.NewAge
	updatedUser := uh.Storage.Update(intId, *user)
	if updatedUser == nil {
		http.Error(w, "Not found", http.StatusNotFound)
		log.Println("Not found: Пользователь не найден")
		return
	}
	wr := fmt.Sprintf("Возраст %s изменен на %d", user.Name, user.Age)
	log.Println("User Updated.", wr)
	w.Write([]byte("Возраст пользователя успешно обновлен"))
}

func (uh *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	var req requests.Delete
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println("Bad request:", err.Error())
		return
	}
	var user models.User
	for _, u := range uh.Storage.Users {
		if u.ID == req.Target_id {
			user = *u
			break
		}
	}
	if uh.Storage.Delete(req.Target_id) == nil {
		http.Error(w, "Not found", http.StatusNotFound)
		log.Println("Not found: Пользователь не найден")
		return
	}
	log.Printf("User deleted. Name=%s", user.Name)
	w.Write([]byte(fmt.Sprint(user.Name)))
}
