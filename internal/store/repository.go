package store

import "testsocket/internal/model"

type UserRepository interface {
	Find(int) (*model.User, error)
	FindByToken(string) (*model.User, error)
}
