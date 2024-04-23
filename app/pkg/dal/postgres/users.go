package dbdal

import (
	"context"
	"fmt"

	dbClient "plato/app/pkg/client/db"
)

// User represents the structure for a user record in the database.
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

const usersTable = "users"

// GetUserById retrieves a user by their ID using direct SQL queries.
func GetUserById(ctx context.Context, id string) (*User, error) {
	query := fmt.Sprintf("SELECT id, name, email FROM %s WHERE id = $1", usersTable)
	row := dbClient.GetPGClient().QueryRow(ctx, query, id)

	var user User
	if err := row.Scan(&user.ID, &user.Name, &user.Email); err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}
	return &user, nil
}

// AddUser adds a new user to the database with the provided name and email.
func AddUser(ctx context.Context, name, email string) (*User, error) {
	query := fmt.Sprintf("INSERT INTO %s (name, email) VALUES ($1, $2) RETURNING id, name, email", usersTable)
	row := dbClient.GetPGClient().QueryRow(ctx, query, name, email)

	var user User
	if err := row.Scan(&user.ID, &user.Name, &user.Email); err != nil {
		return nil, fmt.Errorf("error adding user: %w", err)
	}
	return &user, nil
}
