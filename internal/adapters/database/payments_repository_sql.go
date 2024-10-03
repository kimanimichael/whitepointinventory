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

type PaymentRepositorySql struct {
	DB *sqlcdatabase.Queries
}

var _ domain.PaymentsRepository = (*PaymentRepositorySql)(nil)

func NewPaymentsRepositorySQL(DB *sqlcdatabase.Queries) *PaymentRepositorySql {
	return &PaymentRepositorySql{DB: DB}
}

func (r *PaymentRepositorySql) CreatePayment(cashPaid, chickenPrice int32, farmerName string, user *domain.User) (*domain.Payment, error) {
	farmer, err := r.DB.GetFarmerByName(context.Background(), farmerName)
	if err != nil {
		return nil, err
	}
	payment, err := r.DB.CreatePayment(context.Background(), sqlcdatabase.CreatePaymentParams{
		ID:                  uuid.New(),
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
		CashPaid:            cashPaid,
		PricePerChickenPaid: chickenPrice,
		UserID:              user.ID,
		FarmerID:            farmer.ID,
	})
	if err != nil {
		return nil, err
	}
	cashBalance := sql.NullInt32{}
	chickenBalance := sql.NullFloat64{}

	cashBalance.Int32 = cashPaid
	cashBalance.Valid = true

	chickenBalance.Float64 = float64(cashPaid) / float64(chickenPrice)

	err = r.DB.DecreaseChickenOwed(context.Background(), sqlcdatabase.DecreaseChickenOwedParams{
		ChickenBalance: chickenBalance,
		ID:             farmer.ID,
	})
	if err != nil {
		_ = fmt.Errorf("couldn't decrease chicken owed: %v", err)
	}
	err = r.DB.DecreaseCashOwed(context.Background(), sqlcdatabase.DecreaseCashOwedParams{
		CashBalance: cashBalance,
		ID:          farmer.ID,
	})
	if err != nil {
		_ = fmt.Errorf("couldn't decrease cash owed: %v", err)
	}

	err = r.DB.MarkFarmerAsUpdated(context.Background(), farmer.ID)
	if err != nil {
		_ = fmt.Errorf("couldn't mark farmer as updated: %v", err)
	}

	fmt.Printf("Farmer %v updated at updated to %v\n", farmer.Name, time.Now())

	modelPayment := models.DatabasePaymentToPayment(payment)
	return &domain.Payment{
		ID:                   modelPayment.ID,
		CreatedAt:            modelPayment.CreatedAt,
		UpdatedAt:            modelPayment.UpdatedAt,
		CashPaid:             modelPayment.CashPaid,
		PricePerChickenPaid:  modelPayment.PricePerChickenPaid,
		FarmerID:             modelPayment.FarmerID,
		UserName:             modelPayment.UserName,
		FarmerName:           modelPayment.FarmerName,
		FarmerChickenBalance: modelPayment.FarmerChickenBalance,
		FarmerCashBalance:    modelPayment.FarmerCashBalance,
	}, nil
}

func (r *PaymentRepositorySql) GetPaymentByID(ID uuid.UUID) (*domain.Payment, error) {
	//TODO implement me
	panic("implement me")
}

func (r *PaymentRepositorySql) GetPayments() ([]domain.Payment, error) {
	//TODO implement me
	panic("implement me")
}

func (r *PaymentRepositorySql) DeletePayment(ID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}
