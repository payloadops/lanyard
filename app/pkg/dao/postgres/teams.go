package postgresdao

type Team struct {
	TeamId   int
	TeamName string
	OrgId    int
}

func (dao *DAO) GetTeamById(teamId int) (*Team, error) {
	team := &Team{}
	err := dao.db.QueryRow("SELECT team_id, team_name, org_id FROM teams WHERE team_id = $1", teamId).Scan(&team.TeamId, &team.TeamName, &team.OrgId)
	if err != nil {
		return nil, err
	}
	return team, nil
}

func (dao *DAO) AddTeam(team *Team) error {
	_, err := dao.db.Exec("INSERT INTO teams (team_name, org_id) VALUES ($1, $2)", team.TeamName, team.OrgId)
	return err
}
