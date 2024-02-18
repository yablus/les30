package dto

import (
	"github.com/yablus/les30/internal/models"
)

func AddIdToUser(r models.ReqCreate) models.User {
	var user models.User
	models.IDs++
	user.ID = models.IDs
	user.Name = r.Name
	user.Age = r.Age
	user.Friends = r.Friends
	return user
}

func ListFriends(users []*models.User, friends []int) string {
	var list string
	for _, u := range users {
		for _, v := range friends {
			if u.ID == v {
				if list != "" {
					list += ", "
				}
				list += u.Name
			}
		}
	}
	return list
}
