package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"testing"

	"github.com/bmcszk/user-service/db"
	"github.com/bmcszk/user-service/logic"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/joho/godotenv"
)

type Block struct {
	*testing.T
	ctx        context.Context
	client     *http.Client
	serviceUri string
	queries    *db.Queries

	givenID   int64
	givenUser logic.User
	request   *http.Request

	response      *http.Response
	returnedUser  *logic.User
	returnedUsers *logic.UsersResponse
	returnErr     error
}

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		slog.Warn("failed to load .env")
	}
}

func NewBlocks(t *testing.T) (*Block, *Block, *Block) {
	postgresUrl := os.Getenv("POSTGRES_URL")
	ctx, cancel := context.WithCancel(context.Background())
	conn, err := pgx.Connect(ctx, postgresUrl)
	if err != nil {
		t.Fatal(err)
	}
	queries := db.New(conn)
	t.Cleanup(func() {
		// TODO clean DB after tests
		cancel()
		conn.Close(ctx)
	})
	b := &Block{
		T:          t,
		ctx:        ctx,
		client:     http.DefaultClient,
		queries:    queries,
		serviceUri: os.Getenv("SERVICE_URI"),
	}
	return b, b, b
}

func (b *Block) and() *Block {
	return b
}

func (b *Block) aValidUserData() *Block {
	b.givenUser = logic.User{
		Name:  randomString(10),
		Other: "e2e test user",
	}
	return b
}

func (b *Block) aID() *Block {
	b.givenID = int64(rand.Intn(10000)) + 10000
	return b
}

func (b *Block) userIsChanged() *Block {
	b.givenUser.Name = randomString(10)
	b.givenUser.Other = "e2e test user changed"
	return b
}

func (b *Block) aInvalidUserData() *Block {
	b.givenUser = logic.User{
		Name:  "",
		Other: "other",
	}
	return b
}

func (b *Block) postRequest() *Block {
	requestBody, err := json.Marshal(b.givenUser)
	if err != nil {
		b.Fatal(err)
	}
	b.request, err = http.NewRequestWithContext(b.ctx, http.MethodPost, fmt.Sprintf("%s/users", b.serviceUri), bytes.NewReader(requestBody))
	if err != nil {
		b.Fatal(err)
	}
	return b
}

func (b *Block) getRequest() *Block {
	var err error
	b.request, err = http.NewRequestWithContext(b.ctx, http.MethodGet, fmt.Sprintf("%s/users/%v", b.serviceUri, b.givenID), nil)
	if err != nil {
		b.Fatal(err)
	}
	return b
}

func (b *Block) putRequest() *Block {
	requestBody, err := json.Marshal(b.givenUser)
	if err != nil {
		b.Fatal(err)
	}
	b.request, err = http.NewRequestWithContext(b.ctx, http.MethodPut, fmt.Sprintf("%s/users/%v", b.serviceUri, b.givenID), bytes.NewReader(requestBody))
	if err != nil {
		b.Fatal(err)
	}
	return b
}

func (b *Block) deleteRequest() *Block {
	var err error
	b.request, err = http.NewRequestWithContext(b.ctx, http.MethodDelete, fmt.Sprintf("%s/users/%v", b.serviceUri, b.givenID), nil)
	if err != nil {
		b.Fatal(err)
	}
	return b
}

func (b *Block) listRequest() *Block {
	var err error
	b.request, err = http.NewRequestWithContext(b.ctx, http.MethodGet, fmt.Sprintf("%s/users", b.serviceUri), nil)
	if err != nil {
		b.Fatal(err)
	}
	return b
}

func (b *Block) sending() *Block {
	b.response, b.returnErr = b.client.Do(b.request)
	return b
}

func (b *Block) noError() *Block {
	if b.returnErr != nil {
		b.Fatal(b.returnErr)
	}
	return b
}

func (b *Block) statusCodeIs(code int) *Block {
	if b.response.StatusCode != code {
		b.Fatalf("status code not expected: %v", b.response.StatusCode)
	}
	return b
}

func (b *Block) userIsReturned() *Block {
	err := json.NewDecoder(b.response.Body).Decode(&b.returnedUser)
	if err != nil {
		b.Fatal(err)
	}
	defer b.response.Body.Close()
	return b
}

func (b *Block) returnedUserIsValid() *Block {
	if b.returnedUser == nil {
		b.Fatal("user not returned")
	}
	if b.returnedUser.ID == 0 || b.returnedUser.CreatedAt.IsZero() {
		b.Fatal("invalid user")
	}
	return b
}

func (b *Block) userIsValid(user *logic.User) *Block {
	if user == nil {
		b.Fatal("user not returned")
	}
	if user.ID == 0 || user.CreatedAt.IsZero() {
		b.Fatal("invalid user")
	}
	return b
}

func (b *Block) usersAreReturned() *Block {
	err := json.NewDecoder(b.response.Body).Decode(&b.returnedUsers)
	if err != nil {
		b.Fatal(err)
	}
	defer b.response.Body.Close()
	return b
}

func (b *Block) returnedUsersAreValid() *Block {
	if b.returnedUsers.Count == 0 {
		b.Fatal("users not returned")
	}
	if b.returnedUsers.Count != len(b.returnedUsers.Users) {
		b.Fatal("invalid number of users returned")
	}
	for _, user := range b.returnedUsers.Users {
		b.userIsValid(user)
	}
	return b
}

func (b *Block) alreadyStoredInDB() *Block {
	dbUser, err := b.queries.CreateUser(b.ctx, db.CreateUserParams{
		Name:  b.givenUser.Name,
		Other: pgtype.Text{String: b.givenUser.Other, Valid: true},
	})
	if err != nil {
		b.Fatal(err)
	}
	b.givenID = dbUser.ID
	return b
}

func (b *Block) userIsStoredInDB() *Block {
	dbUser, err := b.queries.GetUser(b.ctx, b.returnedUser.ID)
	if err != nil {
		b.Fatal(err)
	}
	if dbUser.Name != b.givenUser.Name || dbUser.Other.String != b.givenUser.Other {
		b.Fatal("user stored in db not expected")
	}
	return b
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
