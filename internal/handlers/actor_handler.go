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

func ActorHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "actors handler done")
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
	// Получаем ID актера из URL-адреса
	urlPath := r.URL.Path
	idStr := strings.TrimPrefix(urlPath, "/actors/delete/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid actor ID", http.StatusBadRequest)
		return
	}

	// Удаляем актера из базы данных
	query := `DELETE FROM actors WHERE id=$1`
	_, err = db.Exec(query, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем статус 204 No Content
	w.WriteHeader(http.StatusNoContent)
}
