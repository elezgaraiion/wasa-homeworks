package database

import "github.com/aritz/wasa-homeworks/service/models"

func (db *appdbimpl) GetUserByID(id string) (models.User, error) {
	row := db.c.QueryRow(`
		SELECT id, name, photo FROM users WHERE id = ?
	`, id)

	var u models.User
	err := row.Scan(&u.ID, &u.Name, &u.Photo)
	if err != nil {
		return models.User{}, err
	}
	return u, nil
}