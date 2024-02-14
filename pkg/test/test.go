package test

import "github.com/yablus/les30/internal/models"

var fakeUsers = []*models.User{
	{
		ID:      1,
		Name:    "some name",
		Age:     25,
		Friends: []int{2, 3},
	}, {
		ID:      2,
		Name:    "another name",
		Age:     23,
		Friends: []int{1},
	}, {
		ID:      3,
		Name:    "third name",
		Age:     31,
		Friends: []int{1},
	},
}

type FakeStorage struct {
	Users []*models.User
}

func (s *FakeStorage) List() []*models.User {
	return fakeUsers
}

func (s *FakeStorage) Get(_ int) *models.User {
	return fakeUsers[0]
}

func (s *FakeStorage) Update(int, models.User) *models.User {
	return fakeUsers[1]
}

func (s *FakeStorage) Create(_ models.User) {
	return
}

func (s *FakeStorage) Delete(_ int) *models.User {
	return nil
}
