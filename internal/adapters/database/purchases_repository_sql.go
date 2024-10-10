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
		fmt.Printf("Farmer %v updated at %v\n", farmer.Name, time.Now())
	}

	updatedFarmer, err := r.DB.GetFarmerByName(context.Background(), farmerName)
	if err != nil {
		_ = fmt.Errorf("couldn't get farmer by name: %v", err)
	}

	modelPurchase := models.DatabasePurchaseToPurchase(purchase)
	return &domain.Purchase{
		ID:                   modelPurchase.ID,
		CreatedAt:            modelPurchase.CreatedAt,
		UpdatedAt:            modelPurchase.UpdatedAt,
		FarmerID:             modelPurchase.FarmerID,
		UserID:               modelPurchase.UserID,
		FarmerName:           farmer.Name,
		UserName:             user.Name,
		Chicken:              modelPurchase.Chicken,
		PricePerChicken:      modelPurchase.PricePerChicken,
		FarmerChickenBalance: updatedFarmer.ChickenBalance.Float64,
		FarmerCashBalance:    updatedFarmer.CashBalance.Int32,
	}, nil
}

func (r *PurchasesRepositorySQL) GetPurchaseByID(ID uuid.UUID) (*domain.Purchase, error) {
	purchase, err := r.DB.GetPurchaseByID(context.Background(), ID)
	if err != nil {
		return nil, err
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

func (r *PurchasesRepositorySQL) GetPurchases() ([]domain.Purchase, error) {
	purchases, err := r.DB.GetPurchases(context.Background())
	if err != nil {
		return nil, err
	}
	var purchasesWithNames []models.Purchase
	var purchasesToReturn []models.Purchase
	for _, purchase := range purchases {
		purchasesWithNames = append(purchasesWithNames, models.DatabasePurchaseToPurchase(purchase))
	}
	for _, purchaseWithName := range purchasesWithNames {
		user, err := r.DB.GetUserByID(context.Background(), purchaseWithName.UserID)
		if err != nil {
			return nil, err
		}
		purchaseWithName.UserName = user.Name
		farmer, err := r.DB.GetFarmerByID(context.Background(), purchaseWithName.FarmerID)
		if err != nil {
			return nil, err
		}
		purchaseWithName.FarmerName = farmer.Name
		purchaseWithName.FarmerChickenBalance = farmer.ChickenBalance.Float64
		purchaseWithName.FarmerCashBalance = farmer.CashBalance.Int32

		purchasesToReturn = append(purchasesToReturn, purchaseWithName)
	}

	var domainPurchases []domain.Purchase
	for _, purchase := range purchasesToReturn {
		domainPurchases = append(domainPurchases, domain.Purchase{
			ID:                   purchase.ID,
			CreatedAt:            purchase.CreatedAt,
			UpdatedAt:            purchase.UpdatedAt,
			Chicken:              purchase.Chicken,
			PricePerChicken:      purchase.PricePerChicken,
			UserID:               purchase.UserID,
			FarmerID:             purchase.FarmerID,
			UserName:             purchase.UserName,
			FarmerName:           purchase.FarmerName,
			FarmerChickenBalance: purchase.FarmerChickenBalance,
			FarmerCashBalance:    purchase.FarmerCashBalance,
		})
	}
	return domainPurchases, nil
}

func (r *PurchasesRepositorySQL) DeletePurchase(ID uuid.UUID) error {
	cashBalance := sql.NullInt32{}
	chickenBalance := sql.NullFloat64{}

	purchase, err := r.DB.GetPurchaseByID(context.Background(), ID)
	if err != nil {
		return err
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
		_ = fmt.Errorf("couldn't increase cash owed: %v", err)
		return err
	}
	err = r.DB.DecreaseChickenOwed(context.Background(), sqlcdatabase.DecreaseChickenOwedParams{
		ChickenBalance: chickenBalance,
		ID:             purchase.FarmerID,
	})
	if err != nil {
		_ = fmt.Errorf("couldn't increase chicken owed: %v", err)
		return err
	}
	err = r.DB.DeletePayments(context.Background(), ID)
	if err != nil {
		_ = fmt.Errorf("couldn't delete purchase: %v", err)
		return err
	}
	return nil
}
