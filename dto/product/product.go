package productdto

type CreateProductRequest struct {
	Title      string `json:"title" validate:"required"`
	Price      int    `json:"price" validate:"required"`
	Image      string `json:"image"`
	CategoryID int    `json:"category_id"`
}

type UpdateProductRequest struct {
	Title      string `json:"title"`
	Price      int    `json:"price"`
	Image      string `json:"image"`
	CategoryID int    `json:"category_id"`
}

type ProductResponse struct {
	ID         int      `json:"id"`
	Title      string   `json:"title"`
	Price      string   `json:"price"`
	Image      string   `json:"image"`
	Admin      string   `json:"admin"`
	CategoryID []int    `json:"category_id"`
	Category   []string `json:"category"`
}
