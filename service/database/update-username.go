package database

import (
	"errors"
	"log"
	"database/sql"
	"github.com/aritz/wasa-homeworks/service/models"
)

// UpdateUserName cambia el nombre de un usuario por su ID.
// Devuelve el usuario actualizado o un error si el nombre ya existe o falla la base de datos.
func (db *appdbimpl) UpdateUserName(id, newName string) (models.User, error) {
	// Verificar si el nuevo nombre ya existe para otro usuario
	var exists int
	err := db.c.QueryRow(`SELECT COUNT(*) FROM users WHERE name = ? AND id != ?`, newName, id).Scan(&exists)
	if err != nil {
		return models.User{}, err
	}
	if exists > 0 {
		return models.User{}, errors.New("user with that name already exists")
	}

	// Actualizar el nombre
	log.Printf("UpdateUserName called with id=%s, newName=%s", id, newName)

	_, err = db.c.Exec(`UPDATE users SET name = ? WHERE id = ?`, newName, id)
	if err != nil {
    	log.Printf("UPDATE error: %v", err)
    	return models.User{}, err
	}

	var u models.User
	err = db.c.QueryRow(`SELECT id, name, photo FROM users WHERE id = ?`, id).Scan(&u.ID, &u.Name, &u.Photo)
	if err == sql.ErrNoRows {
    	log.Printf("User not found after update: %s", id)
    	return models.User{}, errors.New("user not found")
	}
	if err != nil {
    	log.Printf("SELECT after update error: %v", err)
    	return models.User{}, err
	}

	return u, nil
}