package main

import (
	"github.com/7Maliko7/to-do-list-api/internal/api"
	"github.com/7Maliko7/to-do-list-api/internal/storage"
	"log"
	"net/http"
)

func main() {
	store := storage.New("http://localhost:8080")
	api.Store = &store
	http.HandleFunc("/list", api.GetListHandler)
	http.HandleFunc("/create", api.CreateHandler)
	http.HandleFunc("/get", api.GetHandler)
	http.HandleFunc("/update", api.UpdateHandler)
	http.HandleFunc("/delete", api.DeleteHandler)
	log.Fatal(http.ListenAndServe(":8081", nil))
}
