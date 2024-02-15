package repository

import "github.com/yablus/les30/internal/models"

type MemStorage struct {
	Users []*models.User
}

func NewStorage() *MemStorage {
	return &MemStorage{
		Users: make([]*models.User, 0),
	}
}

func (u *MemStorage) List() []*models.User {
	return u.Users
}

func (u *MemStorage) Get(id int) *models.User {
	for _, user := range u.Users {
		if user.ID == id {
			return user
		}
	}
	return nil
}

func (u *MemStorage) Update(id int, userUpdate models.User) *models.User {
	for i, user := range u.Users {
		if user.ID == id {
			u.Users[i] = &userUpdate
			return user
		}
	}
	return nil
}

func (u *MemStorage) Create(user models.User) {
	u.Users = append(u.Users, &user)
}

func (u *MemStorage) Delete(id int) *models.User {
	for _, user := range u.Users {
		for i, v := range user.Friends {
			if v == id {
				user.Friends = append(user.Friends[:i], (user.Friends)[i+1:]...)
			}
		}
	}
	for i, user := range u.Users {
		if user.ID == id {
			u.Users = append(u.Users[:i], (u.Users)[i+1:]...)
			return &models.User{}
		}
	}
	return nil
}
