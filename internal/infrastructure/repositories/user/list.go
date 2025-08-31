package user

import (
	"context"
	"fmt"
	"slices"

	"github.com/sonalys/goshare/internal/domain"
	"github.com/sonalys/goshare/internal/infrastructure/postgres/sqlcgen"
)

func (r *Repository) ListByEmail(ctx context.Context, emails []string) ([]domain.User, error) {
	emails = slices.Compact(emails)
	users, err := r.conn.Queries().ListByEmail(ctx, emails)
	if err != nil {
		return nil, userError(err)
	}

	var errs domain.Form
	for idx, email := range emails {
		if !slices.ContainsFunc(users, func(user sqlcgen.User) bool {
			return user.Email == email
		}) {
			errs = append(errs, domain.FieldError{
				Field: fmt.Sprintf("emails.%d", idx),
				Cause: domain.ErrUserNotFound,
			})
		}
	}

	if err := errs.Close(); err != nil {
		return nil, fmt.Errorf("failed to get users by email: %w", err)
	}

	return toUsers(users), nil
}
