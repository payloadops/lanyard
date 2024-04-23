package dbdal

import (
	"context"
	"fmt"

	dbClient "plato/app/pkg/client/db"
)

// Team represents the structure for a team record in the database.
type Team struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	OrgID string `json:"org_id"` // Assuming there is a column for OrgID in the database schema.
}

const teamsTable = "teams"

// GetTeamById retrieves a team by its ID using direct SQL queries.
func GetTeamById(ctx context.Context, id string) (*Team, error) {
	query := fmt.Sprintf("SELECT id, name, org_id FROM %s WHERE id = $1", teamsTable)
	row := dbClient.GetPGClient().QueryRow(ctx, query, id)

	var team Team
	if err := row.Scan(&team.ID, &team.Name, &team.OrgID); err != nil {
		return nil, fmt.Errorf("error getting team: %w", err)
	}
	return &team, nil
}

// AddTeam adds a new team to the database with the provided name and organization ID.
func AddTeam(ctx context.Context, name, orgId string) (*Team, error) {
	query := fmt.Sprintf("INSERT INTO %s (name, org_id) VALUES ($1, $2) RETURNING id, name, org_id", teamsTable)
	row := dbClient.GetPGClient().QueryRow(ctx, query, name, orgId)

	var team Team
	if err := row.Scan(&team.ID, &team.Name, &team.OrgID); err != nil {
		return nil, fmt.Errorf("error adding team: %w", err)
	}
	return &team, nil
}
