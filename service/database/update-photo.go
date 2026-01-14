package database

import "github.com/aritz/wasa-homeworks/service/models"

func (db *appdbimpl) UpdateMyPhoto(id, photoURL string) (models.User, error) {
	_, err := db.c.Exec(`UPDATE users SET photo = ? WHERE id = ?`, photoURL, id)
	if err != nil {
		return models.User{}, err
	}

	var u models.User
	err = db.c.QueryRow(`SELECT id, name, photo FROM users WHERE id = ?`, id).Scan(&u.ID, &u.Name, &u.Photo)
	if err != nil {
		return models.User{}, err
	}

	return u, nil
}
