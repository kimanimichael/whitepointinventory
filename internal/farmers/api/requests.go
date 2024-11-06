package farmersapi

type CreateFarmerRequest struct {
	Name           string  `json:"name"`
	ChickenBalance float64 `json:"chicken_balance"`
	CashBalance    int32   `json:"cash_balance"`
}

type GetFarmerRequest struct {
	Name string `json:"name"`
}
