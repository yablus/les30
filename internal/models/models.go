package models

var IDs int

type User struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Friends []int  `json:"friends"`
}

type Storage struct {
	Users []*User
}

func NewStorage() *Storage {
	return &Storage{
		Users: make([]*User, 0),
	}
}

func (u *Storage) List() []*User {
	return u.Users
}

func (u *Storage) Get(id int) *User {
	for _, user := range u.Users {
		if user.ID == id {
			return user
		}
	}
	return nil
}

func (u *Storage) Update(id int, userUpdate User) *User {
	for i, user := range u.Users {
		if user.ID == id {
			u.Users[i] = &userUpdate
			return user
		}
	}
	return nil
}

func (u *Storage) Create(user User) {
	u.Users = append(u.Users, &user)
}

func (u *Storage) Delete(id int) *User {
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
			return &User{}
		}
	}
	return nil
}
