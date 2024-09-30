package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	sqlcdatabase "github.com/mike-kimani/whitepointinventory/internal/adapters/database/sqlc/gensql"
	"github.com/mike-kimani/whitepointinventory/internal/domain"
	"github.com/mike-kimani/whitepointinventory/internal/models"
	"time"
)

type PurchasesRepositorySQL struct {
	DB *sqlcdatabase.Queries
}

var _ domain.PurchasesRepository = (*PurchasesRepositorySQL)(nil)

func NewPurchasesRepositorySQL(db *sqlcdatabase.Queries) *PurchasesRepositorySQL {
	return &PurchasesRepositorySQL{
		DB: db,
	}
}

func (r *PurchasesRepositorySQL) CreatePurchase(chickenNo, chickenPrice int32, farmerName string, user *domain.User) (*domain.Purchase, error) {
	farmer, err := r.DB.GetFarmerByName(context.Background(), farmerName)
	if err != nil {
		return nil, err
	}
	purchase, err := r.DB.CreatePurchase(context.Background(), sqlcdatabase.CreatePurchaseParams{
		ID:              uuid.New(),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		Chicken:         chickenNo,
		PricePerChicken: chickenPrice,
		UserID:          user.ID,
		FarmerID:        farmer.ID,
	})
	if err != nil {
		return nil, err
	}

	var chickenBought sql.NullFloat64
	chickenBought.Float64 = float64(chickenNo)
	chickenBought.Valid = true
	err = r.DB.IncreaseChickenOwed(context.Background(), sqlcdatabase.IncreaseChickenOwedParams{
		ChickenBalance: chickenBought,
		ID:             farmer.ID,
	})
	if err != nil {
		return nil, err
	}

	var cashBalance sql.NullInt32
	cashBalance.Int32 = chickenNo * chickenPrice
	cashBalance.Valid = true
	err = r.DB.IncreaseCashOwed(context.Background(), sqlcdatabase.IncreaseCashOwedParams{
		CashBalance: cashBalance,
		ID:          farmer.ID,
	})
	if err != nil {
		return nil, err
	}

	if err = r.DB.MarkFarmerAsUpdated(context.Background(), farmer.ID); err != nil {
		fmt.Printf("Could not mark farmer as updated: %v\n", err)
	}
	if err == nil {
		fmt.Printf("Farmer %v updated at updated to %v\n", farmer.Name, time.Now())
	}

	modelPurchase := models.DatabasePurchaseToPurchase(purchase)
	return &domain.Purchase{
		ID:              modelPurchase.ID,
		CreatedAt:       modelPurchase.CreatedAt,
		UpdatedAt:       modelPurchase.UpdatedAt,
		FarmerID:        modelPurchase.FarmerID,
		UserID:          modelPurchase.UserID,
		Chicken:         modelPurchase.Chicken,
		PricePerChicken: modelPurchase.PricePerChicken,
	}, nil
}
