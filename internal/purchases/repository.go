package purchases

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	sqlcdatabase "github.com/mike-kimani/whitepointinventory/internal/adapters/database/sqlc/gensql"
	"github.com/mike-kimani/whitepointinventory/internal/users"
	"time"
)

type PurchaseRepositorySQL struct {
	DB *sqlcdatabase.Queries
}

var _ PurchaseRepository = (*PurchaseRepositorySQL)(nil)

func NewPurchaseRepositorySQL(db *sqlcdatabase.Queries) *PurchaseRepositorySQL {
	return &PurchaseRepositorySQL{
		DB: db,
	}
}

func (r *PurchaseRepositorySQL) CreatePurchase(chickenNo, chickenPrice int32, farmerName string, user *users.User) (*Purchase, error) {
	farmer, err := r.DB.GetFarmerByName(context.Background(), farmerName)
	if err != nil {
		return nil, fmt.Errorf("could not get farmer: %w", err)
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
		return nil, fmt.Errorf("could not create purchase: %w", err)
	}

	var chickenBought sql.NullFloat64
	chickenBought.Float64 = float64(chickenNo)
	chickenBought.Valid = true
	err = r.DB.IncreaseChickenOwed(context.Background(), sqlcdatabase.IncreaseChickenOwedParams{
		ChickenBalance: chickenBought,
		ID:             farmer.ID,
	})
	if err != nil {
		_ = fmt.Errorf("couldn't decrease chicken owed: %v", err)
	}

	var cashBalance sql.NullInt32
	cashBalance.Int32 = chickenNo * chickenPrice
	cashBalance.Valid = true
	err = r.DB.IncreaseCashOwed(context.Background(), sqlcdatabase.IncreaseCashOwedParams{
		CashBalance: cashBalance,
		ID:          farmer.ID,
	})
	if err != nil {
		_ = fmt.Errorf("couldn't decrease chicken owed: %v", err)
	}

	if err = r.DB.MarkFarmerAsUpdated(context.Background(), farmer.ID); err != nil {
		fmt.Printf("Could not mark farmer as updated: %v\n", err)
	}
	if err == nil {
		fmt.Printf("Farmer %v updated at %v\n", farmer.Name, time.Now())
	}

	updatedFarmer, err := r.DB.GetFarmerByName(context.Background(), farmerName)
	if err != nil {
		_ = fmt.Errorf("couldn't get farmer by name: %v", err)
	}

	return &Purchase{
		ID:                   purchase.ID,
		CreatedAt:            purchase.CreatedAt,
		UpdatedAt:            purchase.UpdatedAt,
		FarmerID:             purchase.FarmerID,
		UserID:               purchase.UserID,
		FarmerName:           farmer.Name,
		UserName:             user.Name,
		Chicken:              purchase.Chicken,
		PricePerChicken:      purchase.PricePerChicken,
		FarmerChickenBalance: updatedFarmer.ChickenBalance.Float64,
		FarmerCashBalance:    updatedFarmer.CashBalance.Int32,
	}, nil
}

func (r *PurchaseRepositorySQL) GetPurchaseByID(ID uuid.UUID) (*Purchase, error) {
	purchase, err := r.DB.GetPurchaseByID(context.Background(), ID)
	if err != nil {
		return nil, fmt.Errorf("error getting purchase ID: %v", err)
	}
	return &Purchase{
		ID:              purchase.ID,
		CreatedAt:       purchase.CreatedAt,
		UpdatedAt:       purchase.UpdatedAt,
		FarmerID:        purchase.FarmerID,
		UserID:          purchase.UserID,
		Chicken:         purchase.Chicken,
		PricePerChicken: purchase.PricePerChicken,
	}, nil
}

func (r *PurchaseRepositorySQL) GetMostRecentPurchase() (*Purchase, error) {
	purchase, err := r.DB.GetMostRecentPurchase(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error getting most recent purchase: %v", err)
	}

	farmer, err := r.DB.GetFarmerByID(context.Background(), purchase.FarmerID)
	if err != nil {
		return nil, fmt.Errorf("error getting farmer from most recent purchase: %v", err)
	}
	return &Purchase{
		ID:              purchase.ID,
		CreatedAt:       purchase.CreatedAt,
		UpdatedAt:       purchase.UpdatedAt,
		FarmerID:        purchase.FarmerID,
		FarmerName:      farmer.Name,
		UserID:          purchase.UserID,
		Chicken:         purchase.Chicken,
		PricePerChicken: purchase.PricePerChicken,
	}, nil
}

func (r *PurchaseRepositorySQL) GetPurchases() ([]Purchase, error) {
	purchases, err := r.DB.GetPurchases(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error getting purchases: %v", err)
	}

	var purchasesToReturn []Purchase

	for _, purchase := range purchases {
		purchaseWithName := &Purchase{
			ID:              purchase.ID,
			CreatedAt:       purchase.CreatedAt,
			UpdatedAt:       purchase.UpdatedAt,
			Chicken:         purchase.Chicken,
			PricePerChicken: purchase.PricePerChicken,
			FarmerID:        purchase.FarmerID,
		}
		user, err := r.DB.GetUserByID(context.Background(), purchase.UserID)
		if err != nil {
			return nil, fmt.Errorf("error getting user from purchase: %v", err)
		}
		purchaseWithName.UserName = user.Name
		farmer, err := r.DB.GetFarmerByID(context.Background(), purchaseWithName.FarmerID)
		if err != nil {
			return nil, fmt.Errorf("error getting farmer from purchase: %v", err)
		}
		purchaseWithName.FarmerName = farmer.Name
		purchaseWithName.FarmerChickenBalance = farmer.ChickenBalance.Float64
		purchaseWithName.FarmerCashBalance = farmer.CashBalance.Int32

		purchasesToReturn = append(purchasesToReturn, *purchaseWithName)
	}

	return purchasesToReturn, nil
}

func (r *PurchaseRepositorySQL) DeletePurchase(ID uuid.UUID) error {
	cashBalance := sql.NullInt32{}
	chickenBalance := sql.NullFloat64{}

	purchase, err := r.DB.GetPurchaseByID(context.Background(), ID)
	if err != nil {
		return fmt.Errorf("error getting purchase from ID: %v", err)
	}
	chickenBalance.Float64 = float64(purchase.Chicken)
	chickenBalance.Valid = true
	cashBalance.Int32 = purchase.Chicken * purchase.PricePerChicken
	cashBalance.Valid = true

	err = r.DB.DecreaseCashOwed(context.Background(), sqlcdatabase.DecreaseCashOwedParams{
		CashBalance: cashBalance,
		ID:          purchase.FarmerID,
	})
	if err != nil {
		return fmt.Errorf("couldn't decrease cash owed: %v", err)
	}
	err = r.DB.DecreaseChickenOwed(context.Background(), sqlcdatabase.DecreaseChickenOwedParams{
		ChickenBalance: chickenBalance,
		ID:             purchase.FarmerID,
	})
	if err != nil {
		return fmt.Errorf("couldn't decrease chicken owed: %v", err)
	}
	err = r.DB.DeletePurchase(context.Background(), ID)
	if err != nil {
		return fmt.Errorf("couldn't delete purchase: %v", err)
	}
	return nil
}
