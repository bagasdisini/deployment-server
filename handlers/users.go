package handlers

import (
	dto "backend/dto/result"
	usersdto "backend/dto/users"
	"backend/models"
	"backend/repositories"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

type handler struct {
	UserRepository repositories.UserRepository
}

func HandlerUser(UserRepository repositories.UserRepository) *handler {
	return &handler{UserRepository}
}

func (h *handler) ShowUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	users, err := h.UserRepository.ShowUsers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
	}

	for i, p := range users {
		imagePath := os.Getenv("PATH_FILE") + p.Image
		users[i].Image = imagePath
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Status: http.StatusOK, Data: users}
	json.NewEncoder(w).Encode(response)
}

func (h *handler) ShowAdmins(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	users, err := h.UserRepository.ShowAdmins()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
	}

	for i, p := range users {
		imagePath := os.Getenv("PATH_FILE") + p.Image
		users[i].Image = imagePath
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Status: http.StatusOK, Data: users}
	json.NewEncoder(w).Encode(response)
}

func (h *handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	user, err := h.UserRepository.GetUserByIDUser(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	if user.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		response := dto.ErrorResult{Status: http.StatusNotFound, Message: "ID: " + strconv.Itoa(id) + " not found!"}
		json.NewEncoder(w).Encode(response)
		return
	}

	user.Image = os.Getenv("PATH_FILE") + user.Image

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Status: http.StatusOK, Data: user}
	json.NewEncoder(w).Encode(response)
}

func (h *handler) GetAdminByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	user, err := h.UserRepository.GetUserByIDAdmin(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	if user.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		response := dto.ErrorResult{Status: http.StatusNotFound, Message: "ID: " + strconv.Itoa(id) + " not found!"}
		json.NewEncoder(w).Encode(response)
		return
	}

	user.Image = os.Getenv("PATH_FILE") + user.Image

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Status: http.StatusOK, Data: user}
	json.NewEncoder(w).Encode(response)
}

func (h *handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	dataUpload := r.Context().Value("dataFile")
	filepath := ""

	var ctx = context.Background()
	var CLOUD_NAME = os.Getenv("CLOUD_NAME")
	var API_KEY = os.Getenv("API_KEY")
	var API_SECRET = os.Getenv("API_SECRET")

	if dataUpload != nil {
		filepath = dataUpload.(string)
	}

	cld, _ := cloudinary.NewFromParams(CLOUD_NAME, API_KEY, API_SECRET)

	resp, err2 := cld.Upload.Upload(ctx, filepath, uploader.UploadParams{Folder: "dumbmerch"})

	if err2 != nil {
		fmt.Println(err2.Error())
	}

	request := usersdto.UpdateUserRequest{
		FullName: r.FormValue("fullName"),
		Email:    r.FormValue("email"),
		Phone:    r.FormValue("phone"),
		Location: r.FormValue("location"),
	}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	user := models.User{
		Image: resp.SecureURL,
	}

	user.ID = id

	if request.FullName != "" {
		user.FullName = request.FullName
	}

	if request.Email != "" {
		user.Email = request.Email
	}

	if request.Phone != "" {
		user.Phone = request.Phone
	}

	if filepath != "" {
		user.Image = filepath
	}

	if request.Location != "" {
		user.Location = request.Location
	}

	if request.Gender != "" {
		user.Gender = request.Gender
	}

	data, err := h.UserRepository.UpdateUser(user, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Status: http.StatusOK, Data: data}
	json.NewEncoder(w).Encode(response)
}

func (h *handler) UpdateAdmin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	dataUpload := r.Context().Value("dataFile")
	filename := ""

	if dataUpload != nil {
		filename = dataUpload.(string)
	}

	request := usersdto.UpdateUserRequest{
		FullName: r.FormValue("fullName"),
		Email:    r.FormValue("email"),
		Phone:    r.FormValue("phone"),
		Location: r.FormValue("location"),
	}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	user := models.Admin{}

	user.ID = id

	if request.FullName != "" {
		user.FullName = request.FullName
	}

	if request.Email != "" {
		user.Email = request.Email
	}

	if filename != "" {
		user.Image = filename
	}

	if request.Phone != "" {
		user.Phone = request.Phone
	}

	if request.Location != "" {
		user.Location = request.Location
	}

	if request.Gender != "" {
		user.Gender = request.Gender
	}

	data, err := h.UserRepository.UpdateAdmin(user, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Status: http.StatusOK, Data: data}
	json.NewEncoder(w).Encode(response)
}

func (h *handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	roleInfo := r.Context().Value("authInfo").(jwt.MapClaims)
	roleStr := string(roleInfo["role"].(string))

	if roleStr == "user" {
		user, err := h.UserRepository.GetUserByIDUser(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			response := dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()}
			json.NewEncoder(w).Encode(response)
			return
		}

		data, err := h.UserRepository.DeleteUser(user, id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()}
			json.NewEncoder(w).Encode(response)
			return
		}

		w.WriteHeader(http.StatusOK)
		response := dto.SuccessResult{Status: http.StatusOK, Data: data}
		json.NewEncoder(w).Encode(response)
	}

	if roleStr == "admin" {
		user, err := h.UserRepository.GetUserByIDAdmin(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			response := dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()}
			json.NewEncoder(w).Encode(response)
			return
		}

		data, err := h.UserRepository.DeleteAdmin(user, id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()}
			json.NewEncoder(w).Encode(response)
			return
		}

		w.WriteHeader(http.StatusOK)
		response := dto.SuccessResult{Status: http.StatusOK, Data: data}
		json.NewEncoder(w).Encode(response)
	}
}
