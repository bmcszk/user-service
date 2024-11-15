package logic

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bmcszk/user-service/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

var createdAt = time.Now()
var updatedAt = time.Now()

type Block struct {
	*testing.T
	queries *MockQueries
	service *Service

	givenID   int64
	givenUser User

	returnedUser  *User
	returnedUsers *UsersResponse
	returnErr     error
}

func NewBlocks(t *testing.T) (*Block, *Block, *Block) {
	queries := &MockQueries{}
	service := NewService(queries)
	b := &Block{
		T:       t,
		queries: queries,
		service: service,
	}
	return b, b, b
}

func (b *Block) and() *Block {
	return b
}

func (b *Block) aUser() *Block {
	b.givenUser = User{
		Name:  "name",
		Other: "other",
	}
	return b
}

func (b *Block) aInvaliUser() *Block {
	b.givenUser = User{
		Name:  "",
		Other: "other",
	}
	return b
}

func (b *Block) aID() *Block {
	b.givenID = 7
	return b
}

func (b *Block) dbCanCreateUser() *Block {
	b.queries.createUser = func(ctx context.Context, params db.CreateUserParams) (db.User, error) {
		return db.User{
			ID:        1,
			Name:      params.Name,
			Other:     params.Other,
			CreatedAt: pgtype.Timestamp{Time: createdAt, Valid: true},
		}, nil
	}
	return b
}

func (b *Block) dbCannotCreateDuplicatedUser() *Block {
	b.queries.createUser = func(ctx context.Context, params db.CreateUserParams) (db.User, error) {
		return db.User{}, &pgconn.PgError{
			Code:    "23505",
			Message: `duplicate key value violates unique constraint "users_name"`,
		}
	}
	return b
}

func (b *Block) dbCanGetUser() *Block {
	b.queries.getUser = func(ctx context.Context, id int64) (db.User, error) {
		return db.User{
			ID:        id,
			Name:      "name",
			Other:     pgtype.Text{String: "other", Valid: true},
			CreatedAt: pgtype.Timestamp{Time: createdAt, Valid: true},
		}, nil
	}
	return b
}

func (b *Block) dbCannotFindUser() *Block {
	b.queries.getUser = func(ctx context.Context, id int64) (db.User, error) {
		return db.User{}, pgx.ErrNoRows
	}
	return b
}

func (b *Block) dbCanUpdateUser() *Block {
	b.queries.updateUser = func(ctx context.Context, params db.UpdateUserParams) (db.User, error) {
		return db.User{
			ID:        params.ID,
			Name:      params.Name,
			Other:     params.Other,
			CreatedAt: pgtype.Timestamp{Time: createdAt, Valid: true},
			UpdatedAt: pgtype.Timestamp{Time: updatedAt, Valid: true},
		}, nil
	}
	return b
}

func (b *Block) dbCannotUpdateWithExistingUsername() *Block {
	b.queries.updateUser = func(ctx context.Context, params db.UpdateUserParams) (db.User, error) {
		return db.User{}, &pgconn.PgError{
			Code:    "23505",
			Message: `duplicate key value violates unique constraint "users_name"`,
		}
	}
	return b
}

func (b *Block) dbCannotFindUserForUpdate() *Block {
	b.queries.updateUser = func(ctx context.Context, params db.UpdateUserParams) (db.User, error) {
		return db.User{}, pgx.ErrNoRows
	}
	return b
}

func (b *Block) dbCanDeleteUser() *Block {
	b.queries.deleteUser = func(ctx context.Context, id int64) error {
		return nil
	}
	return b
}

func (b *Block) dbCannotFindUserForDelete() *Block {
	b.queries.deleteUser = func(ctx context.Context, id int64) error {
		return pgx.ErrNoRows
	}
	return b
}

func (b *Block) dbCanListUsers() *Block {
	b.queries.listUsers = func(ctx context.Context, params db.ListUsersParams) ([]db.User, error) {
		return []db.User{
			{
				ID:        1,
				Name:      "name",
				Other:     pgtype.Text{String: "other", Valid: true},
				CreatedAt: pgtype.Timestamp{Time: createdAt, Valid: true},
			},
			{
				ID:        2,
				Name:      "name2",
				Other:     pgtype.Text{String: "other2", Valid: true},
				CreatedAt: pgtype.Timestamp{Time: createdAt, Valid: true},
			},
		}, nil
	}
	return b
}

func (b *Block) serviceCreatesUser() *Block {
	b.returnedUser, b.returnErr = b.service.CreateUser(context.Background(), b.givenUser)
	return b
}

func (b *Block) serviceGetsUser() *Block {
	b.returnedUser, b.returnErr = b.service.GetUserByID(context.Background(), b.givenID)
	return b
}

func (b *Block) serviceUpdatesUser() *Block {
	b.returnedUser, b.returnErr = b.service.UpdateUserByID(context.Background(), b.givenID, b.givenUser)
	return b
}

func (b *Block) serviceDeletesUser() *Block {
	err := b.service.DeleteUserByID(context.Background(), b.givenID)
	b.returnErr = err
	return b
}

func (b *Block) serviceListsUsers() *Block {
	b.returnedUsers, b.returnErr = b.service.ListUsers(context.Background(), 10, 0)
	return b
}

func (b *Block) noError() *Block {
	if b.returnErr != nil {
		b.Fatal(b.returnErr)
	}
	return b
}

func (b *Block) returnedErrorIs(err error) *Block {
	if b.returnErr == nil {
		b.Fatal("error not returned")
	}
	if !errors.Is(b.returnErr, err) {
		b.Fatal("error not expected")
	}
	return b
}

func (b *Block) userIsReturned() *Block {
	if b.returnedUser == nil {
		b.Fatal("user not returned")
	}
	return b
}

func (b *Block) returnedUserIsValid() *Block {
	return b.userIsValid(b.returnedUser)
}

func (b *Block) returnedUsersAreValid() *Block {
	if b.returnedUsers.Count != len(b.returnedUsers.Users) {
		b.Fatal("invalid number of users returned")
	}
	for _, user := range b.returnedUsers.Users {
		b.userIsValid(user)
	}
	return b
}

func (b *Block) userIsValid(user *User) *Block {
	if user.ID == 0 || user.CreatedAt.IsZero() {
		b.Fatal("invalid user returned")
	}
	return b
}

func (b *Block) usersAreReturned() *Block {
	if b.returnedUsers == nil || len(b.returnedUsers.Users) == 0 {
		b.Fatal("users not returned")
	}
	return b
}

type MockQueries struct {
	createUser func(context.Context, db.CreateUserParams) (db.User, error)
	getUser    func(context.Context, int64) (db.User, error)
	updateUser func(context.Context, db.UpdateUserParams) (db.User, error)
	deleteUser func(context.Context, int64) error
	listUsers  func(context.Context, db.ListUsersParams) ([]db.User, error)
}

func (m *MockQueries) CreateUser(ctx context.Context, params db.CreateUserParams) (db.User, error) {
	return m.createUser(ctx, params)
}

func (m *MockQueries) GetUser(ctx context.Context, id int64) (db.User, error) {
	return m.getUser(ctx, id)
}

func (m *MockQueries) UpdateUser(ctx context.Context, params db.UpdateUserParams) (db.User, error) {
	return m.updateUser(ctx, params)
}

func (m *MockQueries) DeleteUser(ctx context.Context, id int64) error {
	return m.deleteUser(ctx, id)
}

func (m *MockQueries) ListUsers(ctx context.Context, params db.ListUsersParams) ([]db.User, error) {
	return m.listUsers(ctx, params)
}
