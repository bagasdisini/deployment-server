package transactiondto

type CreateTransactionRequest struct {
	Value   int    `json:"value"`
	Product string `json:"product"`
	BuyerID int    `json:"user_id"`
	AdminID int    `json:"admin_id"`
}

type UpdateTransactionRequest struct {
	Value   int    `json:"value"`
	Product string `json:"product"`
	BuyerID int    `json:"user_id"`
	AdminID int    `json:"admin_id"`
}

type TransactionResponse struct {
	ID      int    `json:"id"`
	Value   int    `json:"value"`
	Status  string `json:"status"`
	BuyerID int    `json:"user_id"`
	AdminID int    `json:"admin_id"`
	Product string `json:"product"`
	Date    string `json:"date"`
}
