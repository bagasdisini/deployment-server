package handlers

import (
	dto "backend/dto/result"
	"backend/repositories"
	"encoding/json"
	"net/http"
)

type handlerCategories struct {
	CategoriesRepository repositories.CategoryRepository
}

func HandlerCategories(CategoriesRepository repositories.CategoryRepository) *handlerCategories {
	return &handlerCategories{CategoriesRepository}
}

func (h *handlerCategories) ShowCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	categories, err := h.CategoriesRepository.ShowCategory()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Status: http.StatusOK, Data: categories}
	json.NewEncoder(w).Encode(response)
}
