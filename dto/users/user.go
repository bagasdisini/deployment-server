package usersdto

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
	FullName string `json:"fullName" validate:"required"`
	Gender   string `json:"gender" validate:"required"`
	Phone    string `json:"phone" validate:"required"`
	Role     string `json:"role" validate:"required"`
}

type UpdateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"fullName"`
	Gender   string `json:"gender"`
	Phone    string `json:"phone"`
	Location string `json:"location"`
	Image    string `json:"image"`
}

type UserResponse struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"fullName"`
	Gender   string `json:"gender"`
	Phone    string `json:"phone"`
	Location string `json:"location"`
	Role     string `json:"role"`
	Image    string `json:"image"`
}
