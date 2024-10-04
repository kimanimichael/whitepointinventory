package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/mike-kimani/whitepointinventory/internal/adapters/database/sqlc/gensql"
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
	payment, err := r.DB.GetPaymentByID(context.Background(), ID)
	if err != nil {
		return nil, err
	}
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

func (r *PaymentRepositorySql) GetPayments() ([]domain.Payment, error) {
	var customPayments []models.Payment
	var paymentResponse []domain.Payment

	payments, err := r.DB.GetPayments(context.Background())
	if err != nil {
		return nil, err
	}
	for _, payment := range payments {
		customPayments = append(customPayments, models.DatabasePaymentToPayment(payment))
	}
	for _, customPayment := range customPayments {
		user, err := r.DB.GetUserByID(context.Background(), customPayment.UserID)
		if err != nil {
			return nil, err
		}
		farmer, err := r.DB.GetFarmerByID(context.Background(), customPayment.FarmerID)
		if err != nil {
			return nil, err
		}
		paymentResponse = append(paymentResponse, domain.Payment{
			ID:                   customPayment.ID,
			CreatedAt:            customPayment.CreatedAt,
			UpdatedAt:            customPayment.UpdatedAt,
			CashPaid:             customPayment.CashPaid,
			PricePerChickenPaid:  customPayment.PricePerChickenPaid,
			FarmerID:             customPayment.FarmerID,
			UserName:             user.Name,
			FarmerName:           farmer.Name,
			FarmerChickenBalance: farmer.ChickenBalance.Float64,
			FarmerCashBalance:    farmer.CashBalance.Int32,
		})
	}
	return paymentResponse, nil
}

func (r *PaymentRepositorySql) DeletePayment(ID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}
