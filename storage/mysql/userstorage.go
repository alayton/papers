package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/alayton/papers"
	"github.com/alayton/papers/storage/mysql/gen"
)

type UserStorage struct {
	db *sql.DB
	q  *gen.Queries
}

func NewUserStorage(db *sql.DB) *UserStorage {
	return &UserStorage{
		db: db,
		q:  gen.New(db),
	}
}

func NewUser() papers.User {
	return &gen.User{}
}

func (s *UserStorage) CreateUser(ctx context.Context, user papers.User) error {
	u, ok := user.(*gen.User)
	if !ok {
		return fmt.Errorf("UserStorage.CreateUser was given the wrong User implementation")
	}

	result, err := s.q.CreateUser(ctx, gen.CreateUserParams{
		Email:         u.Email,
		Password:      u.Password,
		TotpSecret:    u.TotpSecret,
		Confirmed:     u.Confirmed,
		ConfirmToken:  u.ConfirmToken,
		RecoveryToken: u.RecoveryToken,
		LockedUntil:   u.LockedUntil,
		Attempts:      u.Attempts,
		LastAttempt:   u.LastAttempt,
	})
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.SetID(id)

	return nil
}

func (s *UserStorage) UpdateUser(ctx context.Context, user papers.User) error {
	u, ok := user.(*gen.User)
	if !ok {
		return fmt.Errorf("UserStorage.UpdateUser was given the wrong User implementation")
	}

	return s.q.UpdateUser(ctx, gen.UpdateUserParams{
		Email:         u.Email,
		Password:      u.Password,
		TotpSecret:    u.TotpSecret,
		Confirmed:     u.Confirmed,
		ConfirmToken:  u.ConfirmToken,
		RecoveryToken: u.RecoveryToken,
		LockedUntil:   u.LockedUntil,
		Attempts:      u.Attempts,
		LastAttempt:   u.LastAttempt,
		ID:            u.ID,
	})
}

func (s *UserStorage) CreateOAuth2Identity(ctx context.Context, user papers.User, provider, identity string) error {
	return s.q.CreateOAuth2Identity(ctx, gen.CreateOAuth2IdentityParams{
		UserID:   int32(user.GetID()),
		Provider: provider,
		Identity: identity,
	})
}

func (s *UserStorage) RemoveOAuth2Identity(ctx context.Context, user papers.User, provider, identity string) error {
	return s.q.RemoveOAuth2Identity(ctx, gen.RemoveOAuth2IdentityParams{
		UserID:   int32(user.GetID()),
		Provider: provider,
		Identity: identity,
	})
}

func (s *UserStorage) GetUserByID(ctx context.Context, id int64) (papers.User, error) {
	user, err := s.q.GetUser(ctx, id)
	if err == sql.ErrNoRows {
		return nil, papers.ErrUserNotFound
	} else if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserStorage) GetUserByEmail(ctx context.Context, email string) (papers.User, error) {
	user, err := s.q.GetUserByEmail(ctx, email)
	if err == sql.ErrNoRows {
		return nil, papers.ErrUserNotFound
	} else if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserStorage) GetUserByUsername(ctx context.Context, username string) (papers.User, error) {
	user, err := s.q.GetUserByUsername(ctx, sql.NullString{String: username, Valid: true})
	if err == sql.ErrNoRows {
		return nil, papers.ErrUserNotFound
	} else if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserStorage) GetUserByOAuth2Identity(ctx context.Context, provider, identity string) (papers.User, error) {
	user, err := s.q.GetUserByOAuth2Identity(ctx, gen.GetUserByOAuth2IdentityParams{
		Provider: provider,
		Identity: identity,
	})
	if err == sql.ErrNoRows {
		return nil, papers.ErrUserNotFound
	} else if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserStorage) GetUserByConfirmationToken(ctx context.Context, token string) (papers.User, error) {
	user, err := s.q.GetUserByConfirmationToken(ctx, sql.NullString{String: token, Valid: true})
	if err == sql.ErrNoRows {
		return nil, papers.ErrUserNotFound
	} else if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserStorage) GetUserByRecoveryToken(ctx context.Context, token string) (papers.User, error) {
	user, err := s.q.GetUserByRecoveryToken(ctx, sql.NullString{String: token, Valid: true})
	if err == sql.ErrNoRows {
		return nil, papers.ErrUserNotFound
	} else if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserStorage) GetUserRoles(ctx context.Context, user papers.User) ([]string, error) {
	return s.q.GetUserRoles(ctx, int32(user.GetID()))
}

func (s *UserStorage) GetUserPermissions(ctx context.Context, user papers.User) ([]string, error) {
	return s.q.GetUserPermissions(ctx, int32(user.GetID()))
}
