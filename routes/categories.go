package routes

import (
	"backend/handlers"
	"backend/pkg/mysql"
	"backend/repositories"

	"github.com/gorilla/mux"
)

func CategoryRoute(r *mux.Router) {
	CategoryRepository := repositories.RepositoryCategory(mysql.DB)
	h := handlers.HandlerCategories(CategoryRepository)

	r.HandleFunc("/categories", h.ShowCategories).Methods("GET")
}
