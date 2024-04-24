package dbdal

import (
	"context"
	"fmt"
	dbClient "plato/app/pkg/client/db"

	"github.com/uptrace/bun"
)

// Team represents the structure for a team record in the database.
type Team struct {
	bun.BaseModel `bun:"table:teams,alias:t"`
	ID            string `bun:"id,pk" json:"id"`
	Name          string `bun:"name" json:"name"`
	OrgID         string `bun:"org_id" json:"org_id"` // Assuming there is a column for OrgID in the database schema.
}

// GetTeamById retrieves a team by its ID.
func GetTeamById(ctx context.Context, id string) (*Team, error) {
	team := &Team{}
	err := dbClient.GetClient().NewSelect().Model(team).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting team: %w", err)
	}
	return team, nil
}

// AddTeam adds a new team to the database with the provided name and organization ID.
func AddTeam(ctx context.Context, name, orgId string) (*Team, error) {
	team := &Team{
		Name:  name,
		OrgID: orgId,
	}
	_, err := dbClient.GetClient().NewInsert().Model(team).Returning("*").Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("error adding team: %w", err)
	}
	return team, nil
}
