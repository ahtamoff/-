package main

import (
	"film-library/internal/handlers"
	"film-library/internal/storage"
	"log"
	"net/http"
)

func main() {

	db, err := storage.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("/actors/add", func(w http.ResponseWriter, r *http.Request) {
		handlers.AddActorHandler(db, w, r)
	})
	mux.HandleFunc("/actors/update/{id}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			handlers.UpdateActorHandler(db, w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	// Обработчик для удаления актера
	http.HandleFunc("/actors/delete/{id}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			// Здесь вызывайте функцию для удаления актера
			handlers.DeleteActorHandler(db, w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.ListenAndServe(":1247", mux)
}
