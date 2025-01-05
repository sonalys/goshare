package postgres

import (
	"context"

	"github.com/sonalys/goshare/internal/infrastructure/postgres/sqlc"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
)

func (r *ExpensesRepository) CreatePayment(ctx context.Context, payment *v1.ExpensePayment) error {
	return mapError(r.client.queries().CreateExpensePayment(ctx, sqlc.CreateExpensePaymentParams{
		ID:          convertUUID(payment.ID),
		UserID:      convertUUID(payment.PaidByID),
		ExpenseID:   convertUUID(payment.ExpenseID),
		PaymentDate: convertTime(payment.PaymentDate),
		Amount:      payment.Amount,
		CreatedAt:   convertTime(payment.CreatedAt),
		CreatedBy:   convertUUID(payment.CreatedBy),
		UpdatedAt:   convertTime(payment.UpdatedAt),
		UpdatedBy:   convertUUID(payment.UpdatedBy),
	}))
}

func (r *ExpensesRepository) UpdatePayment(ctx context.Context, payment *v1.ExpensePayment) error {
	return mapError(r.client.queries().UpdateExpensePayment(ctx, sqlc.UpdateExpensePaymentParams{
		ID:          convertUUID(payment.ID),
		UserID:      convertUUID(payment.PaidByID),
		PaymentDate: convertTime(payment.PaymentDate),
		Amount:      payment.Amount,
		UpdatedAt:   convertTime(payment.UpdatedAt),
		UpdatedBy:   convertUUID(payment.UpdatedBy),
	}))
}

func (r *ExpensesRepository) DeletePayment(ctx context.Context, id v1.ID) error {
	return mapError(r.client.queries().DeleteExpensePayment(ctx, convertUUID(id)))
}

func (r *ExpensesRepository) GetPayments(ctx context.Context, expenseID v1.ID) ([]v1.ExpensePayment, error) {
	payments, err := r.client.queries().GetExpensePayments(ctx, convertUUID(expenseID))
	if err != nil {
		return nil, mapError(err)
	}

	result := make([]v1.ExpensePayment, 0, len(payments))
	for _, payment := range payments {
		result = append(result, *newExpensePayment(&payment))
	}

	return result, nil
}

func newExpensePayment(payment *sqlc.ExpensePayment) *v1.ExpensePayment {
	return &v1.ExpensePayment{
		ID:          newUUID(payment.ID),
		ExpenseID:   newUUID(payment.ExpenseID),
		LedgerID:    newUUID(payment.LedgerID),
		PaidByID:    newUUID(payment.UserID),
		Amount:      payment.Amount,
		PaymentDate: payment.PaymentDate.Time,
		CreatedAt:   payment.CreatedAt.Time,
		CreatedBy:   newUUID(payment.CreatedBy),
		UpdatedAt:   payment.UpdatedAt.Time,
		UpdatedBy:   newUUID(payment.UpdatedBy),
	}
}
