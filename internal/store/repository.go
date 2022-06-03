package store

import "github.com/eldar34/goWebSocket/internal/model"

type UserRepository interface {
	Find(int) (*model.User, error)
	FindByToken(string) (*model.User, error)
}
