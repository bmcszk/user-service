package logic

import (
	"time"

	"github.com/bmcszk/user-service/db"
)

type User struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	Other     string     `json:"other"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

func FromDBUser(dbUser db.User) *User {
	var updatedAt *time.Time
	if dbUser.UpdatedAt.Valid {
		updatedAt = &dbUser.UpdatedAt.Time
	}
	return &User{
		ID:        dbUser.ID,
		Name:      dbUser.Name,
		Other:     dbUser.Other.String,
		CreatedAt: dbUser.CreatedAt.Time,
		UpdatedAt: updatedAt,
	}
}

type UsersResponse struct {
	Users []*User `json:"users"`
	Count int     `json:"count"`
}

func FromDBUsers(dbUser []db.User) *UsersResponse {
	users := make([]*User, len(dbUser))
	for i, dbUser := range dbUser {
		users[i] = FromDBUser(dbUser)
	}
	return &UsersResponse{
		Users: users,
		Count: len(dbUser),
	}
}
