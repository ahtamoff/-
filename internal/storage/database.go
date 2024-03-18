package storage

import (
    "database/sql"
    _ "github.com/lib/pq"
)

func ConnectDB() (*sql.DB, error) {
    connStr := "user=postgres password=postgres dbname=postgres sslmode=disable"
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, err
    }
    if err = db.Ping(); err != nil {
        return nil, err
    }
    return db, nil
}

func DeleteActor(db *sql.DB, actorID int) error {
    _, err := db.Exec("DELETE FROM actors WHERE id = ?", actorID)
    if err != nil {
        return err
    }
    return nil
}