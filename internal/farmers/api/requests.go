package farmersapi

type CreateFarmerRequest struct {
	Name           string  `json:"name"`
	ChickenBalance float64 `json:"chicken_balance"`
	CashBalance    int32   `json:"cash_balance"`
}

type GetFarmerRequest struct {
	Name string `json:"name"`
}

type GetPagedFarmersRequest struct {
	Offset uint32 `json:"offset"`
	Limit  uint32 `json:"limit"`
}

type SetFarmerBalancesRequest struct {
	Name           string  `json:"name"`
	ChickenBalance float64 `json:"chicken_balance"`
	CashBalance    int32   `json:"cash_balance"`
}
