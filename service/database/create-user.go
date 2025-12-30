package database

import (
	"github.com/aritz/wasa-homeworks/service/models"
)
func (db *appdbimpl) CreateUser(u models.User) error {
    _, err := db.c.Exec(
        "INSERT INTO users(id, name, photo) VALUES (?, ?, ?)",
        u.ID, u.Name, u.Photo,
    )
    return err
}
