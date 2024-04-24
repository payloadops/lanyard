package dbdal

import (
	"context"
	"fmt"
	dbClient "plato/app/pkg/client/db"

	"github.com/uptrace/bun"
)

// Organization represents the structure of an organization record in the database.
type Organization struct {
	bun.BaseModel `bun:"table:organizations,alias:o"`
	ID            string `bun:"id,pk" json:"id"`
	Name          string `bun:"name" json:"name"`
}

// GetOrganizationById retrieves an organization by its ID using Bun.
func GetOrganizationById(ctx context.Context, id string) (*Organization, error) {
	organization := &Organization{}
	err := dbClient.GetClient().NewSelect().Model(organization).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting organization: %w", err)
	}
	return organization, nil
}

// AddOrganization adds a new organization to the database with the provided name.
func AddOrganization(ctx context.Context, name string) (*Organization, error) {
	organization := &Organization{
		Name: name,
	}
	_, err := dbClient.GetClient().NewInsert().Model(organization).Returning("*").Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("error adding organization: %w", err)
	}
	return organization, nil
}
