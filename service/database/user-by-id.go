package database

import (
	"database/sql"
	"errors"

	"github.com/aritz/wasa-homeworks/service/models"
)

func (db *appdbimpl) GetUserByID(id string) (models.User, error) {
	var u models.User
	err := db.c.QueryRow(`
        SELECT id, name, photo FROM users WHERE id = ?
    `, id).Scan(&u.ID, &u.Name, &u.Photo)

	if errors.Is(err, sql.ErrNoRows) {
		return models.User{}, models.ErrUserNotFound
	}
	if err != nil {
		return models.User{}, err
	}
	return u, nil
}