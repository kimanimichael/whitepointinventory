package payments

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/mike-kimani/whitepointinventory/internal/adapters/database/sqlc/gensql"
	"github.com/mike-kimani/whitepointinventory/internal/models"
	"github.com/mike-kimani/whitepointinventory/internal/users"
	"time"
)

type PaymentRepositorySql struct {
	DB *sqlcdatabase.Queries
}

var _ PaymentsRepository = (*PaymentRepositorySql)(nil)

func NewPaymentsRepositorySQL(DB *sqlcdatabase.Queries) *PaymentRepositorySql {
	return &PaymentRepositorySql{DB: DB}
}

func (r *PaymentRepositorySql) CreatePayment(cashPaid, chickenPrice int32, farmerName string, user *users.User) (*Payment, error) {
	farmer, err := r.DB.GetFarmerByName(context.Background(), farmerName)
	if err != nil {
		return nil, fmt.Errorf("error getting farmer: %v", err)
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
		return nil, fmt.Errorf("error creating payment: %v", err)
	}
	cashBalance := sql.NullInt32{}
	chickenBalance := sql.NullFloat64{}

	cashBalance.Int32 = cashPaid
	cashBalance.Valid = true

	chickenBalance.Float64 = float64(cashPaid) / float64(chickenPrice)
	chickenBalance.Valid = true

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
	updatedFarmer, err := r.DB.GetFarmerByName(context.Background(), farmerName)
	if err != nil {
		_ = fmt.Errorf("couldn't get farmer by name: %v", err)
	}

	fmt.Printf("Farmer %v updated at %v\n", farmer.Name, time.Now())

	modelPayment := models.DatabasePaymentToPayment(payment)
	return &Payment{
		ID:                   modelPayment.ID,
		CreatedAt:            modelPayment.CreatedAt,
		UpdatedAt:            modelPayment.UpdatedAt,
		CashPaid:             modelPayment.CashPaid,
		PricePerChickenPaid:  modelPayment.PricePerChickenPaid,
		UserID:               modelPayment.UserID,
		FarmerID:             modelPayment.FarmerID,
		UserName:             user.Name,
		FarmerName:           farmer.Name,
		FarmerChickenBalance: updatedFarmer.ChickenBalance.Float64,
		FarmerCashBalance:    updatedFarmer.CashBalance.Int32,
	}, nil
}

func (r *PaymentRepositorySql) GetPaymentByID(ID uuid.UUID) (*Payment, error) {
	payment, err := r.DB.GetPaymentByID(context.Background(), ID)
	if err != nil {
		return nil, fmt.Errorf("error getting payment from ID: %v", err)
	}
	modelPayment := models.DatabasePaymentToPayment(payment)
	return &Payment{
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

func (r *PaymentRepositorySql) GetMostRecentPayment() (*Payment, error) {
	payment, err := r.DB.GetMostRecentPayment(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error getting most recent payment: %v", err)
	}
	modelPayment := models.DatabasePaymentToPayment(payment)
	farmer, err := r.DB.GetFarmerByID(context.Background(), payment.FarmerID)
	if err != nil {
		return nil, fmt.Errorf("error getting farmer from most recent purchase: %v", err)
	}
	return &Payment{
		ID:                   modelPayment.ID,
		CreatedAt:            modelPayment.CreatedAt,
		UpdatedAt:            modelPayment.UpdatedAt,
		CashPaid:             modelPayment.CashPaid,
		PricePerChickenPaid:  modelPayment.PricePerChickenPaid,
		FarmerID:             modelPayment.FarmerID,
		UserName:             modelPayment.UserName,
		FarmerName:           farmer.Name,
		FarmerChickenBalance: modelPayment.FarmerChickenBalance,
		FarmerCashBalance:    modelPayment.FarmerCashBalance,
	}, nil

}

func (r *PaymentRepositorySql) GetPayments() ([]Payment, error) {
	var customPayments []models.Payment
	var paymentResponse []Payment

	payments, err := r.DB.GetPayments(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error getting payments: %v", err)
	}
	for _, payment := range payments {
		customPayments = append(customPayments, models.DatabasePaymentToPayment(payment))
	}
	for _, customPayment := range customPayments {
		user, err := r.DB.GetUserByID(context.Background(), customPayment.UserID)
		if err != nil {
			return nil, fmt.Errorf("error getting user from payment: %v", err)
		}
		farmer, err := r.DB.GetFarmerByID(context.Background(), customPayment.FarmerID)
		if err != nil {
			return nil, fmt.Errorf("error getting farmer from payment: %v", err)
		}
		paymentResponse = append(paymentResponse, Payment{
			ID:                   customPayment.ID,
			CreatedAt:            customPayment.CreatedAt,
			UpdatedAt:            customPayment.UpdatedAt,
			CashPaid:             customPayment.CashPaid,
			PricePerChickenPaid:  customPayment.PricePerChickenPaid,
			FarmerID:             customPayment.FarmerID,
			UserID:               customPayment.UserID,
			UserName:             user.Name,
			FarmerName:           farmer.Name,
			FarmerChickenBalance: farmer.ChickenBalance.Float64,
			FarmerCashBalance:    farmer.CashBalance.Int32,
		})
	}
	return paymentResponse, nil
}

func (r *PaymentRepositorySql) DeletePayment(ID uuid.UUID) error {
	cashBalance := sql.NullInt32{}
	chickenBalance := sql.NullFloat64{}
	payment, err := r.DB.GetPaymentByID(context.Background(), ID)
	if err != nil {
		return fmt.Errorf("error getting payment from ID: %v", err)
	}
	cashBalance.Int32 = payment.CashPaid
	cashBalance.Valid = true
	chickenBalance.Float64 = float64(payment.CashPaid) / float64(payment.PricePerChickenPaid)
	chickenBalance.Valid = true

	err = r.DB.IncreaseCashOwed(context.Background(), sqlcdatabase.IncreaseCashOwedParams{
		CashBalance: cashBalance,
		ID:          payment.FarmerID,
	})
	if err != nil {
		return fmt.Errorf("error increasing cash owed: %v", err)
	}
	err = r.DB.IncreaseChickenOwed(context.Background(), sqlcdatabase.IncreaseChickenOwedParams{
		ChickenBalance: chickenBalance,
		ID:             payment.FarmerID,
	})
	if err != nil {
		return fmt.Errorf("error increasing chicken owed: %v", err)
	}
	err = r.DB.DeletePayments(context.Background(), ID)
	if err != nil {
		return fmt.Errorf("error deleting payment: %v", err)
	}
	return nil
}
