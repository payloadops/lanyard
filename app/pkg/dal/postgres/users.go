package dbdal

import (
	"context"
	"fmt"
	dbClient "plato/app/pkg/client/db"

	"github.com/uptrace/bun"
)

// User represents the structure for a user record in the database.
type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`
	ID            string `bun:"id,pk" json:"id"`
	Name          string `bun:"name" json:"name"`
	Email         string `bun:"email" json:"email"`
}

// GetUserById retrieves a user by their ID using Bun.
func GetUserById(ctx context.Context, id string) (*User, error) {
	user := &User{}
	err := dbClient.GetClient().NewSelect().Model(user).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}
	return user, nil
}

// AddUser adds a new user to the database with the provided name and email.
func AddUser(ctx context.Context, name, email string) (*User, error) {
	user := &User{
		Name:  name,
		Email: email,
	}
	_, err := dbClient.GetClient().NewInsert().Model(user).Returning("*").Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("error adding user: %w", err)
	}
	return user, nil
}
