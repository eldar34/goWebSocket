package mysqlStore

import (
	"database/sql"
	"testsocket/internal/model"
	"testsocket/internal/store"
)

type UserRepository struct {
	store *Store
}

func (r *UserRepository) Find(id int) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRow(
		"SELECT id, email, encrypted_password FROM oauth_access_tokens WHERE id = $1",
		id,
	).Scan(
		&u.ID,
		&u.User_id,
		&u.Access_token,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return u, nil
}

func (r *UserRepository) FindByToken(token string) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRow(
		"SELECT id, user_id, access_token FROM oauth_access_tokens WHERE access_token = ?",
		token,
	).Scan(
		&u.ID,
		&u.User_id,
		&u.Access_token,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return u, nil
}
