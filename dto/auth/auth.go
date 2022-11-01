package authdto

type AuthRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
	Role     string `json:"role" validate:"required"`
}

type AuthResponse struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"fullName"`
	Token    string `json:"token"`
	Role     string `json:"role"`
	Gender   string `json:"gender"`
	Image    string `json:"image"`
	Phone    string `json:"phone"`
}

type CheckAuthResponse struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"fullName"`
	Token    string `json:"token"`
	Gender   string `json:"gender"`
	Image    string `json:"image"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
}
