package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Actor struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Gender    string `json:"gender"`
	Birthdate string `json:"birthdate"`
}


func AddActorHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var actor Actor

	err := json.NewDecoder(r.Body).Decode(&actor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `INSERT INTO actors (name, gender, birth_date) VALUES ($1, $2, $3) RETURNING id`
	err = db.QueryRow(query, actor.Name, actor.Gender, actor.Birthdate).Scan(&actor.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(actor)
}

func UpdateActorHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	
	// Получаем ID актера из URL-адреса
	urlPath := r.URL.Path
	idStr := strings.TrimPrefix(urlPath, "/actors/update/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid actor ID", http.StatusBadRequest)
		return
	}

	// Создаем пустую структуру для хранения обновленной информации об актере
	var actor Actor

	// Декодируем JSON из тела запроса в структуру
	err = json.NewDecoder(r.Body).Decode(&actor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Обновляем информацию об актере в базе данных
	query := `UPDATE actors SET name=$1, gender=$2, birth_date=$3 WHERE id=$4 RETURNING id, name, gender, birth_date`
	err = db.QueryRow(query, actor.Name, actor.Gender, actor.Birthdate, id).Scan(&actor.ID, &actor.Name, &actor.Gender, &actor.Birthdate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем обновленную информацию об актере в формате JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(actor)
}

func DeleteActorHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    // Извлекаем id из URL
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

    // Удаляем актера из базы данных
    _, err = db.Exec("DELETE FROM actors WHERE id = ?", actorID)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error deleting actor: %v", err), http.StatusInternalServerError)
        return
    }

    // Отправляем успешный статус
    w.WriteHeader(http.StatusOK)
}

func GetActorsHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    rows, err := db.Query("SELECT id, name, gender, birth_date FROM actors")
    if err != nil {
        http.Error(w, fmt.Sprintf("Error fetching actors: %v", err), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var actors []Actor
    for rows.Next() {
        var actor Actor
        if err := rows.Scan(&actor.ID, &actor.Name, &actor.Gender, &actor.Birthdate); err != nil {
            http.Error(w, fmt.Sprintf("Error scanning actor row: %v", err), http.StatusInternalServerError)
            return
        }
        actors = append(actors, actor)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(actors)
}
