package handlers

import (
	productdto "backend/dto/product"
	dto "backend/dto/result"
	"backend/models"
	"backend/repositories"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"context"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

type handlerProduct struct {
	ProductRepository repositories.ProductRepository
}

func HandlerProduct(ProductRepository repositories.ProductRepository) *handlerProduct {
	return &handlerProduct{ProductRepository}
}

func (h *handlerProduct) ShowProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	products, err := h.ProductRepository.ShowProducts()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	for i, p := range products {
		imagePath := os.Getenv("PATH_FILE") + p.Image
		products[i].Image = imagePath
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Status: http.StatusOK, Data: products}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerProduct) GetProductByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var product models.Product
	product, err := h.ProductRepository.GetProductByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	product.Image = os.Getenv("PATH_FILE") + product.Image

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Status: http.StatusOK, Data: product}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerProduct) CreateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	adminInfo := r.Context().Value("authInfo").(jwt.MapClaims)
	adminId := int(adminInfo["id"].(float64))
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

	resp, err2 := cld.Upload.Upload(ctx, filepath, uploader.UploadParams{Folder: "waysfood"})

	if err2 != nil {
		fmt.Println(err2.Error())
	}

	price, _ := strconv.Atoi(r.FormValue("price"))
	category_id, _ := strconv.Atoi(r.FormValue("category_id"))

	request := productdto.CreateProductRequest{
		Title:      r.FormValue("title"),
		Price:      price,
		CategoryID: category_id,
	}

	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	product := models.Product{
		Title:      request.Title,
		Price:      request.Price,
		Image:      resp.SecureURL,
		AdminID:    adminId,
		CategoryID: request.CategoryID,
	}

	product, err = h.ProductRepository.CreateProduct(product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	product, err = h.ProductRepository.GetProductByID(product.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Status: http.StatusOK, Data: product}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerProduct) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	dataUpload := r.Context().Value("dataFile")
	filename := ""

	if dataUpload != nil {
		filename = dataUpload.(string)
	}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	price, _ := strconv.Atoi(r.FormValue("price"))
	category_id, _ := strconv.Atoi(r.FormValue("category_id"))

	request := productdto.UpdateProductRequest{
		Title:      r.FormValue("title"),
		Price:      price,
		CategoryID: category_id,
	}

	product := models.Product{}

	product.ID = id

	if request.Title != "" {
		product.Title = request.Title
	}

	if request.Price != 0 {
		product.Price = request.Price
	}

	if filename != "" {
		product.Image = filename
	}

	if request.CategoryID != 0 {
		product.CategoryID = request.CategoryID
	}

	data, err := h.ProductRepository.UpdateProduct(product, id)
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

func (h *handlerProduct) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	product, err := h.ProductRepository.GetProductByID(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	data, err := h.ProductRepository.DeleteProduct(product, id)
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
