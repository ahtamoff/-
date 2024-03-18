package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Film struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	ReleaseDate string  `json:"release_date"`
	Rating      float64 `json:"rating"`
}


// GetFilmsForActorHandler возвращает список фильмов для указанного актера
func GetFilmsForActorHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    // Извлекаем ID актера из URL
    parts := strings.Split(r.URL.Path, "/")
    if len(parts) < 4 {
        http.Error(w, "Invalid URL", http.StatusBadRequest)
        return
    }

    actorIDStr := parts[len(parts)-1]
    actorID, err := strconv.Atoi(actorIDStr)
    if err != nil {
        http.Error(w, "Invalid actor ID", http.StatusBadRequest)
        return
    }

    rows, err := db.Query("SELECT Films.id, title, description, release_date, rating FROM Films JOIN Film_actors ON Films.id = Film_actors.Film_id WHERE actor_id = $1", actorID)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error fetching Films for actor: %v", err), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var Films []Film
    for rows.Next() {
        var Film Film
        if err := rows.Scan(&Film.ID, &Film.Title, &Film.Description, &Film.ReleaseDate, &Film.Rating); err != nil {
            http.Error(w, fmt.Sprintf("Error scanning Film row: %v", err), http.StatusInternalServerError)
            return
        }
        Films = append(Films, Film)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(Films)
}

// AddFilmHandler добавляет информацию о новом фильме
func AddFilmHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    var Film Film
    err := json.NewDecoder(r.Body).Decode(&Film)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    query := `INSERT INTO Films (title, description, release_date, rating) VALUES ($1, $2, $3, $4) RETURNING id`
    err = db.QueryRow(query, Film.Title, Film.Description, Film.ReleaseDate, Film.Rating).Scan(&Film.ID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(Film)
}

// DeleteFilmHandler удаляет информацию о фильме
func DeleteFilmHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    parts := strings.Split(r.URL.Path, "/")
    if len(parts) < 4 {
        http.Error(w, "Invalid URL", http.StatusBadRequest)
        return
    }

    FilmIDStr := parts[len(parts)-1]
    FilmID, err := strconv.Atoi(FilmIDStr)
    if err != nil {
        http.Error(w, "Invalid Film ID", http.StatusBadRequest)
        return
    }

    _, err = db.Exec("DELETE FROM Films WHERE id = ?", FilmID)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error deleting Film: %v", err), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}

// SearchFilmHandler выполняет поиск фильма по фрагменту названия
func SearchFilmHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query().Get("q")
    if query == "" {
        http.Error(w, "Search query not provided", http.StatusBadRequest)
        return
    }

    rows, err := db.Query("SELECT id, title, description, release_date, rating FROM Films WHERE title LIKE '%' || $1 || '%'", query)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error searching for Films: %v", err), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var Films []Film
    for rows.Next() {
        var Film Film
        if err := rows.Scan(&Film.ID, &Film.Title, &Film.Description, &Film.ReleaseDate, &Film.Rating); err != nil {
            http.Error(w, fmt.Sprintf("Error scanning Film row: %v", err), http.StatusInternalServerError)
            return
        }
        Films = append(Films, Film)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(Films)
}