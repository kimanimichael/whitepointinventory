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

func (r *PurchaseRepositorySQL) CreatePurchase(ctx context.Context, chickenNo, chickenPrice int32, farmerName string, user *users.User) (*Purchase, error) {
	farmer, err := r.DB.GetFarmerByName(ctx, farmerName)
	if err != nil {
		return nil, fmt.Errorf("could not get farmer: %w", err)
	}
	purchase, err := r.DB.CreatePurchase(ctx, sqlcdatabase.CreatePurchaseParams{
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
	err = r.DB.IncreaseChickenOwed(ctx, sqlcdatabase.IncreaseChickenOwedParams{
		ChickenBalance: chickenBought,
		ID:             farmer.ID,
	})
	if err != nil {
		_ = fmt.Errorf("couldn't decrease chicken owed: %v", err)
	}

	var cashBalance sql.NullInt32
	cashBalance.Int32 = chickenNo * chickenPrice
	cashBalance.Valid = true
	err = r.DB.IncreaseCashOwed(ctx, sqlcdatabase.IncreaseCashOwedParams{
		CashBalance: cashBalance,
		ID:          farmer.ID,
	})
	if err != nil {
		_ = fmt.Errorf("couldn't decrease chicken owed: %v", err)
	}

	if err = r.DB.MarkFarmerAsUpdated(ctx, farmer.ID); err != nil {
		fmt.Printf("Could not mark farmer as updated: %v\n", err)
	}
	if err == nil {
		fmt.Printf("Farmer %v updated at %v\n", farmer.Name, time.Now())
	}

	updatedFarmer, err := r.DB.GetFarmerByName(ctx, farmerName)
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

func (r *PurchaseRepositorySQL) GetPurchaseByID(ctx context.Context, ID uuid.UUID) (*Purchase, error) {
	purchase, err := r.DB.GetPurchaseByID(ctx, ID)
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

func (r *PurchaseRepositorySQL) GetMostRecentPurchase(ctx context.Context) (*Purchase, error) {
	purchase, err := r.DB.GetMostRecentPurchase(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting most recent purchase: %v", err)
	}

	farmer, err := r.DB.GetFarmerByID(ctx, purchase.FarmerID)
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

func (r *PurchaseRepositorySQL) GetPurchases(ctx context.Context) ([]Purchase, error) {
	purchases, err := r.DB.GetPurchases(ctx)
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
		user, err := r.DB.GetUserByID(ctx, purchase.UserID)
		if err != nil {
			return nil, fmt.Errorf("error getting user from purchase: %v", err)
		}
		purchaseWithName.UserName = user.Name
		farmer, err := r.DB.GetFarmerByID(ctx, purchaseWithName.FarmerID)
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

func (r *PurchaseRepositorySQL) GetPagedPurchases(ctx context.Context, offset, limit uint32) (*PurchasePage, error) {
	purchases, err := r.DB.GetPagedPurchases(ctx, sqlcdatabase.GetPagedPurchasesParams{
		Offset: int32(offset),
		Limit:  int32(limit),
	})
	if err != nil {
		return nil, fmt.Errorf("error getting paged purchases: %v", err)
	}

	var purchasesToReturn []Purchase
	for _, purchase := range purchases {
		user, err := r.DB.GetUserByID(ctx, purchase.UserID)
		if err != nil {
			return nil, fmt.Errorf("error getting user from purchase: %v", err)
		}
		farmer, err := r.DB.GetFarmerByID(ctx, purchase.FarmerID)
		if err != nil {
			return nil, fmt.Errorf("error getting farmer from purchase: %v", err)
		}
		purchasesToReturn = append(purchasesToReturn, Purchase{
			ID:                   purchase.ID,
			CreatedAt:            purchase.CreatedAt,
			UpdatedAt:            purchase.UpdatedAt,
			Chicken:              purchase.Chicken,
			PricePerChicken:      purchase.PricePerChicken,
			FarmerID:             purchase.FarmerID,
			UserName:             user.Name,
			FarmerName:           farmer.Name,
			FarmerChickenBalance: farmer.ChickenBalance.Float64,
			FarmerCashBalance:    farmer.CashBalance.Int32,
		})
	}

	totalPurchases, err := r.DB.GetPurchasesCount(ctx)
	returnedPage := &PurchasePage{
		page: page{
			Offset: offset,
			Total:  uint32(totalPurchases),
		},
		Purchases: purchasesToReturn,
	}
	return returnedPage, nil
}

func (r *PurchaseRepositorySQL) DeletePurchase(ctx context.Context, ID uuid.UUID) error {
	cashBalance := sql.NullInt32{}
	chickenBalance := sql.NullFloat64{}

	purchase, err := r.DB.GetPurchaseByID(ctx, ID)
	if err != nil {
		return fmt.Errorf("error getting purchase from ID: %v", err)
	}
	chickenBalance.Float64 = float64(purchase.Chicken)
	chickenBalance.Valid = true
	cashBalance.Int32 = purchase.Chicken * purchase.PricePerChicken
	cashBalance.Valid = true

	err = r.DB.DecreaseCashOwed(ctx, sqlcdatabase.DecreaseCashOwedParams{
		CashBalance: cashBalance,
		ID:          purchase.FarmerID,
	})
	if err != nil {
		return fmt.Errorf("couldn't decrease cash owed: %v", err)
	}
	err = r.DB.DecreaseChickenOwed(ctx, sqlcdatabase.DecreaseChickenOwedParams{
		ChickenBalance: chickenBalance,
		ID:             purchase.FarmerID,
	})
	if err != nil {
		return fmt.Errorf("couldn't decrease chicken owed: %v", err)
	}
	err = r.DB.DeletePurchase(ctx, ID)
	if err != nil {
		return fmt.Errorf("couldn't delete purchase: %v", err)
	}
	return nil
}
