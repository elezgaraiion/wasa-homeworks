package database


func (db *appdbimpl) GetUserIdByName(name string) (string, error) {
    var id string
    err := db.c.QueryRow(`SELECT id FROM users WHERE name = ?`, name).Scan(&id)
    if err != nil {
        return "", err
    }
    return id, nil
}
