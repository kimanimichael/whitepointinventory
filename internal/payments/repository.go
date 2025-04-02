package payments

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/mike-kimani/whitepointinventory/internal/adapters/database/sqlc/gensql"
	"github.com/mike-kimani/whitepointinventory/internal/users"
	"log"
	"time"
)

type PaymentRepositorySql struct {
	DB *sqlcdatabase.Queries
}

var _ PaymentsRepository = (*PaymentRepositorySql)(nil)

func NewPaymentsRepositorySQL(DB *sqlcdatabase.Queries) *PaymentRepositorySql {
	return &PaymentRepositorySql{DB: DB}
}

func (r *PaymentRepositorySql) CreatePayment(ctx context.Context, cashPaid, chickenPrice int32, farmerName string, user *users.User) (*Payment, error) {
	farmer, err := r.DB.GetFarmerByName(ctx, farmerName)
	if err != nil {
		return nil, fmt.Errorf("error getting farmer: %v", err)
	}
	userID, err := uuid.Parse(user.ID)
	if err != nil {
		return nil, fmt.Errorf("error parsing user ID: %v", err)
	}
	payment, err := r.DB.CreatePayment(ctx, sqlcdatabase.CreatePaymentParams{
		ID:                  uuid.New(),
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
		CashPaid:            cashPaid,
		PricePerChickenPaid: chickenPrice,
		UserID:              userID,
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

	err = r.DB.DecreaseChickenOwed(ctx, sqlcdatabase.DecreaseChickenOwedParams{
		ChickenBalance: chickenBalance,
		ID:             farmer.ID,
	})
	if err != nil {
		_ = fmt.Errorf("couldn't decrease chicken owed: %v", err)
	}
	err = r.DB.DecreaseCashOwed(ctx, sqlcdatabase.DecreaseCashOwedParams{
		CashBalance: cashBalance,
		ID:          farmer.ID,
	})
	if err != nil {
		_ = fmt.Errorf("couldn't decrease cash owed: %v", err)
	}

	err = r.DB.MarkFarmerAsUpdated(ctx, farmer.ID)
	if err != nil {
		_ = fmt.Errorf("couldn't mark farmer as updated: %v", err)
	}
	updatedFarmer, err := r.DB.GetFarmerByName(ctx, farmerName)
	if err != nil {
		_ = fmt.Errorf("couldn't get farmer by name: %v", err)
	}

	fmt.Printf("Farmer %v updated at %v\n", farmer.Name, time.Now())

	//modelPayment := models.DatabasePaymentToPayment(payment)
	return &Payment{
		ID:                   payment.ID.String(),
		CreatedAt:            payment.CreatedAt,
		UpdatedAt:            payment.UpdatedAt,
		CashPaid:             payment.CashPaid,
		PricePerChickenPaid:  payment.PricePerChickenPaid,
		UserID:               payment.UserID.String(),
		FarmerID:             payment.FarmerID.String(),
		UserName:             user.Name,
		FarmerName:           farmer.Name,
		FarmerChickenBalance: updatedFarmer.ChickenBalance.Float64,
		FarmerCashBalance:    updatedFarmer.CashBalance.Int32,
	}, nil
}

func (r *PaymentRepositorySql) GetPaymentByID(ctx context.Context, ID string) (*Payment, error) {
	uuidID, err := uuid.Parse(ID)
	payment, err := r.DB.GetPaymentByID(ctx, uuidID)
	if err != nil {
		return nil, fmt.Errorf("error getting payment from ID: %v", err)
	}
	//modelPayment := models.DatabasePaymentToPayment(payment)
	return &Payment{
		ID:                  payment.ID.String(),
		CreatedAt:           payment.CreatedAt,
		UpdatedAt:           payment.UpdatedAt,
		CashPaid:            payment.CashPaid,
		PricePerChickenPaid: payment.PricePerChickenPaid,
		FarmerID:            payment.FarmerID.String(),
	}, nil
}

func (r *PaymentRepositorySql) GetMostRecentPayment(ctx context.Context) (*Payment, error) {
	payment, err := r.DB.GetMostRecentPayment(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting most recent payment: %v", err)
	}
	farmer, err := r.DB.GetFarmerByID(ctx, payment.FarmerID)
	if err != nil {
		return nil, fmt.Errorf("error getting farmer from most recent purchase: %v", err)
	}
	return &Payment{
		ID:                   payment.ID.String(),
		CreatedAt:            payment.CreatedAt,
		UpdatedAt:            payment.UpdatedAt,
		CashPaid:             payment.CashPaid,
		PricePerChickenPaid:  payment.PricePerChickenPaid,
		FarmerID:             payment.FarmerID.String(),
		FarmerName:           farmer.Name,
		FarmerChickenBalance: farmer.ChickenBalance.Float64,
		FarmerCashBalance:    farmer.CashBalance.Int32,
	}, nil

}

func (r *PaymentRepositorySql) GetPayments(ctx context.Context) ([]Payment, error) {
	var paymentResponse []Payment

	payments, err := r.DB.GetPayments(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting payments: %v", err)
	}

	for _, payment := range payments {
		user, err := r.DB.GetUserByID(ctx, payment.UserID)
		if err != nil {
			return nil, fmt.Errorf("error getting user from payment: %v", err)
		}
		farmer, err := r.DB.GetFarmerByID(ctx, payment.FarmerID)
		if err != nil {
			return nil, fmt.Errorf("error getting farmer from payment: %v", err)
		}
		paymentResponse = append(paymentResponse, Payment{
			ID:                   payment.ID.String(),
			CreatedAt:            payment.CreatedAt,
			UpdatedAt:            payment.UpdatedAt,
			CashPaid:             payment.CashPaid,
			PricePerChickenPaid:  payment.PricePerChickenPaid,
			FarmerID:             payment.FarmerID.String(),
			UserID:               payment.UserID.String(),
			UserName:             user.Name,
			FarmerName:           farmer.Name,
			FarmerChickenBalance: farmer.ChickenBalance.Float64,
			FarmerCashBalance:    farmer.CashBalance.Int32,
		})
	}
	return paymentResponse, nil
}

func (r *PaymentRepositorySql) GetPagedPayments(ctx context.Context, offset, limit uint32) (*PaymentPage, error) {
	var paymentResponse []Payment
	payments, err := r.DB.GetPagedPayments(ctx, sqlcdatabase.GetPagedPaymentsParams{
		Offset: int32(offset),
		Limit:  int32(limit),
	})
	if err != nil {
		return nil, fmt.Errorf("error getting paged payments: %v", err)
	}
	for _, payment := range payments {
		user, err := r.DB.GetUserByID(ctx, payment.UserID)
		if err != nil {
			return nil, fmt.Errorf("error getting user from payments: %v", err)
		}
		farmer, err := r.DB.GetFarmerByID(ctx, payment.FarmerID)
		if err != nil {
			return nil, fmt.Errorf("error getting farmer from payments: %v", err)
		}
		paymentResponse = append(paymentResponse, Payment{
			ID:                   payment.ID.String(),
			CreatedAt:            payment.CreatedAt,
			UpdatedAt:            payment.UpdatedAt,
			CashPaid:             payment.CashPaid,
			PricePerChickenPaid:  payment.PricePerChickenPaid,
			FarmerID:             payment.FarmerID.String(),
			UserID:               payment.UserID.String(),
			UserName:             user.Name,
			FarmerName:           farmer.Name,
			FarmerChickenBalance: farmer.ChickenBalance.Float64,
			FarmerCashBalance:    farmer.CashBalance.Int32,
		})
	}
	totalPayments, err := r.DB.GetPaymentCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting total payments: %v", err)
	}
	paymentsPage := &PaymentPage{
		page: page{
			Offset: offset,
			Total:  uint32(totalPayments),
		},
		Payments: paymentResponse,
	}
	return paymentsPage, nil
}

func (r *PaymentRepositorySql) ChangePaymentDate(ctx context.Context, paymentID string, date time.Time, user *users.User) error {
	paymentUID, err := uuid.Parse(paymentID)
	if err != nil {
		return fmt.Errorf("error parsing payment id: %v", err)
	}
	err = r.DB.ChangePaymentDate(ctx, sqlcdatabase.ChangePaymentDateParams{
		ID:        paymentUID,
		UpdatedAt: date,
	})
	if err != nil {
		return fmt.Errorf("error updating payment: %v", err)
	}
	log.Printf("payment updated by %s", user.Name)
	return nil
}

func (r *PaymentRepositorySql) DeletePayment(ctx context.Context, ID string) error {
	cashBalance := sql.NullInt32{}
	chickenBalance := sql.NullFloat64{}
	uuidID, err := uuid.Parse(ID)
	payment, err := r.DB.GetPaymentByID(ctx, uuidID)
	if err != nil {
		return fmt.Errorf("error getting payment from ID: %v", err)
	}
	cashBalance.Int32 = payment.CashPaid
	cashBalance.Valid = true
	chickenBalance.Float64 = float64(payment.CashPaid) / float64(payment.PricePerChickenPaid)
	chickenBalance.Valid = true

	err = r.DB.IncreaseCashOwed(ctx, sqlcdatabase.IncreaseCashOwedParams{
		CashBalance: cashBalance,
		ID:          payment.FarmerID,
	})
	if err != nil {
		return fmt.Errorf("error increasing cash owed: %v", err)
	}
	err = r.DB.IncreaseChickenOwed(ctx, sqlcdatabase.IncreaseChickenOwedParams{
		ChickenBalance: chickenBalance,
		ID:             payment.FarmerID,
	})
	if err != nil {
		return fmt.Errorf("error increasing chicken owed: %v", err)
	}
	err = r.DB.DeletePayments(ctx, uuidID)
	if err != nil {
		return fmt.Errorf("error deleting payment: %v", err)
	}
	return nil
}
