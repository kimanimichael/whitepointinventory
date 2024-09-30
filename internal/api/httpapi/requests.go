package httpapi

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email_address"`
	Password string `json:"password"`
}

type CreateFarmerRequest struct {
	Name           string  `json:"name"`
	ChickenBalance float64 `json:"chicken_balance"`
	CashBalance    int32   `json:"cash_balance"`
}

type GetFarmerRequest struct {
	Name string `json:"name"`
}

type CreatePurchaseRequest struct {
	ChickenNo    int32  `json:"chicken"`
	ChickenPrice int32  `json:"chicken_price"`
	FarmerName   string `json:"farmer_name"`
}
