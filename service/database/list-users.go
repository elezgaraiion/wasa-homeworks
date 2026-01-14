package database

import (
	"github.com/aritz/wasa-homeworks/service/models"
)

func (db *appdbimpl) ListUsers(query, currentUserID string) ([]models.User, error) {
	rows, err := db.c.Query(`
        SELECT id, name, photo 
        FROM users
        WHERE name LIKE ? AND id != ?
        ORDER BY name ASC
    `, query+"%", currentUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		err := rows.Scan(&u.ID, &u.Name, &u.Photo)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	if users == nil {
		users = []models.User{}
	}
	return users, nil
}
