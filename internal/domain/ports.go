package domain

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Email     string
	APIKey    string
	Password  string
}

type Farmer struct {
	ID             uuid.UUID
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Name           string
	PhoneNumber    string
	ChickenBalance float64
	CashBalance    int32
}

type Purchase struct {
	ID                   uuid.UUID
	CreatedAt            time.Time
	UpdatedAt            time.Time
	Chicken              int32
	PricePerChicken      int32
	UserID               uuid.UUID
	FarmerID             uuid.UUID
	UserName             string
	FarmerName           string
	FarmerChickenBalance float64
	FarmerCashBalance    int32
}

type Payment struct {
	ID                   uuid.UUID
	CreatedAt            time.Time
	UpdatedAt            time.Time
	CashPaid             int32
	PricePerChickenPaid  int32
	UserID               uuid.UUID
	FarmerID             uuid.UUID
	UserName             string
	FarmerName           string
	FarmerChickenBalance float64
	FarmerCashBalance    int32
}
