package handlers

import (
	authdto "backend/dto/auth"
	dto "backend/dto/result"
	usersdto "backend/dto/users"
	"backend/models"
	"backend/pkg/bcrypt"
	jwtToken "backend/pkg/jwt"
	"backend/repositories"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
)

type handlerAuth struct {
	AuthRepository repositories.AuthRepository
}

func HandlerAuth(AuthRepository repositories.AuthRepository) *handlerAuth {
	return &handlerAuth{AuthRepository}
}

func (h *handlerAuth) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	request := new(usersdto.CreateUserRequest)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	password, err := bcrypt.HashingPassword(request.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
	}

	if request.Role == "user" {
		user := models.User{
			Email:    request.Email,
			Password: password,
			FullName: request.FullName,
			Gender:   request.Gender,
			Phone:    request.Phone,
			Role:     request.Role,
		}

		data, _ := h.AuthRepository.RegisterUser(user)

		w.WriteHeader(http.StatusOK)
		response := dto.SuccessResult{Status: http.StatusOK, Data: data}
		json.NewEncoder(w).Encode(response)
	}

	if request.Role == "admin" {
		admin := models.Admin{
			Email:    request.Email,
			Password: password,
			FullName: request.FullName,
			Gender:   request.Gender,
			Phone:    request.Phone,
			Role:     request.Role,
		}

		data, _ := h.AuthRepository.RegisterAdmin(admin)

		w.WriteHeader(http.StatusOK)
		response := dto.SuccessResult{Status: http.StatusOK, Data: data}
		json.NewEncoder(w).Encode(response)
	}
}

func (h *handlerAuth) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	request := new(authdto.AuthRequest)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	if request.Role == "user" {
		user := models.User{
			Email:    request.Email,
			Password: request.Password,
		}

		user, err := h.AuthRepository.LoginUser(user.Email)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			response := dto.ErrorResult{Status: http.StatusBadRequest, Message: "Email not registered!"}
			json.NewEncoder(w).Encode(response)
			return
		}

		isValid := bcrypt.CheckPasswordHash(request.Password, user.Password)
		if !isValid {
			w.WriteHeader(http.StatusBadRequest)
			response := dto.ErrorResult{Status: http.StatusBadRequest, Message: "Wrong password!"}
			json.NewEncoder(w).Encode(response)
			return
		}

		gnrtToken := jwt.MapClaims{}
		gnrtToken["id"] = user.ID
		gnrtToken["exp"] = time.Now().Add(time.Hour * 3).Unix()
		gnrtToken["role"] = user.Role
		gnrtToken["image"] = user.Image

		token, err := jwtToken.GenerateToken(&gnrtToken)
		if err != nil {
			fmt.Println("Unauthorize")
			return
		}

		AuthResponse := authdto.AuthResponse{
			ID:       user.ID,
			FullName: user.FullName,
			Email:    user.Email,
			Password: user.Password,
			Token:    token,
			Role:     user.Role,
			Gender:   user.Gender,
			Image:    "http://localhost:5000/uploads/" + user.Image,
			Phone:    user.Phone,
		}

		w.Header().Set("Content-Type", "application/json")
		response := dto.SuccessResult{Status: http.StatusOK, Data: AuthResponse}
		json.NewEncoder(w).Encode(response)
	}

	if request.Role == "admin" {
		admin := models.Admin{
			Email:    request.Email,
			Password: request.Password,
		}

		admin, err := h.AuthRepository.LoginAdmin(admin.Email)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			response := dto.ErrorResult{Status: http.StatusBadRequest, Message: "Email not registered!"}
			json.NewEncoder(w).Encode(response)
			return
		}

		isValid := bcrypt.CheckPasswordHash(request.Password, admin.Password)
		if !isValid {
			w.WriteHeader(http.StatusBadRequest)
			response := dto.ErrorResult{Status: http.StatusBadRequest, Message: "Wrong password!"}
			json.NewEncoder(w).Encode(response)
			return
		}

		gnrtToken := jwt.MapClaims{}
		gnrtToken["id"] = admin.ID
		gnrtToken["exp"] = time.Now().Add(time.Hour * 3).Unix()
		gnrtToken["role"] = admin.Role
		gnrtToken["image"] = admin.Image

		token, err := jwtToken.GenerateToken(&gnrtToken)
		if err != nil {
			fmt.Println("Unauthorize")
			return
		}

		AuthResponse := authdto.AuthResponse{
			ID:       admin.ID,
			FullName: admin.FullName,
			Email:    admin.Email,
			Password: admin.Password,
			Token:    token,
			Role:     admin.Role,
			Gender:   admin.Gender,
			Image:    "http://localhost:5000/uploads/" + admin.Image,
			Phone:    admin.Phone,
		}

		w.Header().Set("Content-Type", "application/json")
		response := dto.SuccessResult{Status: http.StatusOK, Data: AuthResponse}
		json.NewEncoder(w).Encode(response)
	}
}

func (h *handlerAuth) CheckAuth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo := r.Context().Value("authInfo").(jwt.MapClaims)
	userId := int(userInfo["id"].(float64))
	userRole := userInfo["role"].(string)

	if userRole == "user" {
		user, err := h.AuthRepository.GetUsers(userId)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			response := dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()}
			json.NewEncoder(w).Encode(response)
			return
		}

		CheckAuthResponse := authdto.CheckAuthResponse{
			ID:       user.ID,
			Email:    user.Email,
			Password: user.Password,
			FullName: user.FullName,
			Role:     user.Role,
			Gender:   user.Gender,
			Phone:    user.Phone,
			Image:    "http://localhost:5000/uploads/" + user.Image,
		}

		w.Header().Set("Content-Type", "application/json")
		response := dto.SuccessResult{Status: http.StatusOK, Data: CheckAuthResponse}
		json.NewEncoder(w).Encode(response)
	}

	if userRole == "admin" {
		user, err := h.AuthRepository.GetAdmins(userId)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			response := dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()}
			json.NewEncoder(w).Encode(response)
			return
		}

		CheckAuthResponse := authdto.CheckAuthResponse{
			ID:       user.ID,
			Email:    user.Email,
			Password: user.Password,
			FullName: user.FullName,
			Role:     user.Role,
			Gender:   user.Gender,
			Phone:    user.Phone,
			Image:    "http://localhost:5000/uploads/" + user.Image,
		}

		w.Header().Set("Content-Type", "application/json")
		response := dto.SuccessResult{Status: http.StatusOK, Data: CheckAuthResponse}
		json.NewEncoder(w).Encode(response)
	}
}
