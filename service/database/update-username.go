package database

import (
	"database/sql"
	"errors"

	"github.com/aritz/wasa-homeworks/service/models"
)

func (db *appdbimpl) UpdateUserName(id, newName string) (models.User, error) {
	var exists int
	err := db.c.QueryRow(`SELECT COUNT(*) FROM users WHERE name = ? AND id != ?`, newName, id).Scan(&exists)
	if err != nil {
		return models.User{}, err
	}
	if exists > 0 {
		return models.User{}, models.ErrNameConflict
	}

	_, err = db.c.Exec(`UPDATE users SET name = ? WHERE id = ?`, newName, id)
	if err != nil {
		return models.User{}, err
	}

	var u models.User
	err = db.c.QueryRow(`SELECT id, name, photo FROM users WHERE id = ?`, id).Scan(&u.ID, &u.Name, &u.Photo)
	if errors.Is(err, sql.ErrNoRows) {
		return models.User{}, models.ErrUserNotFound
	}
	if err != nil {
		return models.User{}, err
	}

	return u, nil
}
