package dbdal

import (
	"context"
	"fmt"

	dbClient "plato/app/pkg/client/db"
)

// Organization represents the structure of an organization record in the database.
type Organization struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

const organizationsTable = "organizations"

// GetOrganizationById retrieves an organization by its ID using direct SQL queries.
func GetOrganizationById(ctx context.Context, id string) (*Organization, error) {
	query := fmt.Sprintf("SELECT id, name FROM %s WHERE id = $1", organizationsTable)
	row := dbClient.GetPGClient().QueryRow(ctx, query, id)

	var organization Organization
	if err := row.Scan(&organization.ID, &organization.Name); err != nil {
		return nil, fmt.Errorf("error getting organization: %w", err)
	}
	return &organization, nil
}

// AddOrganization adds a new organization to the database with the provided name.
func AddOrganization(ctx context.Context, name string) (*Organization, error) {
	query := fmt.Sprintf("INSERT INTO %s (name) VALUES ($1) RETURNING id, name", organizationsTable)
	row := dbClient.GetPGClient().QueryRow(ctx, query, name)

	var organization Organization
	if err := row.Scan(&organization.ID, &organization.Name); err != nil {
		return nil, fmt.Errorf("error adding organization: %w", err)
	}
	return &organization, nil
}
