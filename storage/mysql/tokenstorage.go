package mysql

import (
	"context"
	"database/sql"
	"time"

	"github.com/alayton/papers"
	"github.com/alayton/papers/storage/mysql/gen"
)

type TokenStorage struct {
	db *sql.DB
	q  *gen.Queries
}

func NewTokenStorage(db *sql.DB) *TokenStorage {
	return &TokenStorage{
		db: db,
		q:  gen.New(db),
	}
}

func (s *TokenStorage) CreateAccessToken(ctx context.Context, userID int64, token, chain string, valid bool) error {
	var validInt int32 = 1
	if !valid {
		validInt = 0
	}
	return s.q.CreateAccessToken(ctx, gen.CreateAccessTokenParams{
		UserID:    int32(userID),
		Token:     token,
		Chain:     chain,
		Valid:     validInt,
		CreatedAt: time.Now(),
	})
}

func (s *TokenStorage) CreateRefreshToken(ctx context.Context, userID int64, token, chain string, valid bool) error {
	var validInt int32 = 1
	if !valid {
		validInt = 0
	}
	return s.q.CreateRefreshToken(ctx, gen.CreateRefreshTokenParams{
		UserID:    int32(userID),
		Token:     token,
		Chain:     chain,
		Valid:     validInt,
		CreatedAt: time.Now(),
	})
}

func (s *TokenStorage) GetAccessToken(ctx context.Context, userID int64, token string) (papers.Token, error) {
	t, err := s.q.GetAccessToken(ctx, gen.GetAccessTokenParams{
		UserID: int32(userID),
		Token:  token,
	})
	if err == sql.ErrNoRows {
		return nil, papers.ErrTokenNotFound
	} else if err != nil {
		return nil, err
	}
	return &t, nil
}

func (s *TokenStorage) GetRefreshToken(ctx context.Context, userID int64, token string) (papers.Token, error) {
	t, err := s.q.GetRefreshToken(ctx, gen.GetRefreshTokenParams{
		UserID: int32(userID),
		Token:  token,
	})
	if err == sql.ErrNoRows {
		return nil, papers.ErrTokenNotFound
	} else if err != nil {
		return nil, err
	}
	return &t, nil
}

func (s *TokenStorage) InvalidateAccessTokens(ctx context.Context, userID int64) error {
	return s.q.InvalidateAccessTokens(ctx, int32(userID))
}

func (s *TokenStorage) InvalidateRefreshToken(ctx context.Context, userID int64, token string) error {
	return s.q.InvalidateRefreshToken(ctx, gen.InvalidateRefreshTokenParams{
		UserID: int32(userID),
		Token:  token,
	})
}

func (s *TokenStorage) InvalidateRefreshTokens(ctx context.Context, userID int64) error {
	return s.q.InvalidateRefreshTokens(ctx, int32(userID))
}

func (s *TokenStorage) InvalidateTokenChain(ctx context.Context, userID int64, chain string) error {
	if err := s.q.InvalidateAccessTokenChain(ctx, gen.InvalidateAccessTokenChainParams{
		UserID: int32(userID),
		Chain:  chain,
	}); err != nil {
		return err
	}
	return s.q.InvalidateRefreshTokenChain(ctx, gen.InvalidateRefreshTokenChainParams{
		UserID: int32(userID),
		Chain:  chain,
	})
}

func (s *TokenStorage) PruneAccessTokens(ctx context.Context, timeToStale time.Duration) error {
	stale := time.Now().Add(-timeToStale)
	return s.q.PruneAccessTokens(ctx, stale)
}

func (s *TokenStorage) PruneRefreshTokens(ctx context.Context, timeToStale time.Duration) error {
	stale := time.Now().Add(-timeToStale)
	return s.q.PruneRefreshTokens(ctx, stale)
}
