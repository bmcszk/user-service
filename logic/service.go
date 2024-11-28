package logic

import (
	"context"
	"errors"

	"github.com/bmcszk/user-service/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

const DuplicateErrorCode = "23505"

var ErrUserNotFound = errors.New("user not found")
var ErrUserAlreadyExists = errors.New("user already exists")
var ErrUserNameEmpty = errors.New("user name empty")

type userRepo interface {
	CreateUser(context.Context, db.CreateUserParams) (db.User, error)
	GetUser(context.Context, int64) (db.User, error)
	UpdateUser(context.Context, db.UpdateUserParams) (db.User, error)
	DeleteUser(context.Context, int64) error
	ListUsers(context.Context, db.ListUsersParams) ([]db.User, error)
}

type Service struct {
	userRepo userRepo
}

func NewService(userRepo userRepo) *Service {
	return &Service{
		userRepo: userRepo,
	}
}

func (s *Service) CreateUser(ctx context.Context, user User) (*User, error) {
	if err := validateUser(user); err != nil {
		return nil, err
	}
	dbUser, err := s.userRepo.CreateUser(ctx, db.CreateUserParams{
		Name:  user.Name,
		Other: pgtype.Text{String: user.Other, Valid: true},
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == DuplicateErrorCode {
			return nil, ErrUserAlreadyExists
		}
		return nil, err
	}
	return FromDBUser(dbUser), nil
}

func (s *Service) GetUserByID(ctx context.Context, id int64) (*User, error) {
	dbUser, err := s.userRepo.GetUser(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return FromDBUser(dbUser), nil
}

func (s *Service) UpdateUserByID(ctx context.Context, id int64, user User) (*User, error) {
	if err := validateUser(user); err != nil {
		return nil, err
	}
	dbUser, err := s.userRepo.UpdateUser(ctx, db.UpdateUserParams{
		ID:    id,
		Name:  user.Name,
		Other: pgtype.Text{String: user.Other, Valid: true},
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrUserNotFound
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == DuplicateErrorCode {
			return nil, ErrUserAlreadyExists
		}
		return nil, err
	}
	return FromDBUser(dbUser), nil
}

func (s *Service) DeleteUserByID(ctx context.Context, id int64) error {
	if err := s.userRepo.DeleteUser(ctx, id); err != nil {
		if err == pgx.ErrNoRows {
			return ErrUserNotFound
		}
		return err
	}
	return nil
}

func (s *Service) ListUsers(ctx context.Context, limit, offset int32) (*UsersResponse, error) {
	dbUsers, err := s.userRepo.ListUsers(ctx, db.ListUsersParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}
	return FromDBUsers(dbUsers), nil
}

func validateUser(user User) error {
	if user.Name == "" {
		return ErrUserNameEmpty
	}
	// ...
	return nil
}
