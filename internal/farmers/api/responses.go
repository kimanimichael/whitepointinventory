package farmersapi

import (
	"github.com/google/uuid"
	"github.com/mike-kimani/whitepointinventory/internal/farmers"
	"time"
)

type Farmer struct {
	ID             uuid.UUID `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Name           string    `json:"name"`
	ChickenBalance float64   `json:"chicken_balance"`
	CashBalance    int32     `json:"cash_balance"`
}

func farmerToResponseFarmer(domainFarmer farmers.Farmer) Farmer {
	return Farmer{
		ID:             domainFarmer.ID,
		CreatedAt:      domainFarmer.CreatedAt,
		UpdatedAt:      domainFarmer.UpdatedAt,
		Name:           domainFarmer.Name,
		ChickenBalance: domainFarmer.ChickenBalance,
		CashBalance:    domainFarmer.CashBalance,
	}
}

func farmersToResponseFarmers(domainFarmers []farmers.Farmer) []Farmer {
	var farmers []Farmer
	for _, domainFarmer := range domainFarmers {
		farmers = append(farmers, farmerToResponseFarmer(domainFarmer))
	}
	return farmers
}
