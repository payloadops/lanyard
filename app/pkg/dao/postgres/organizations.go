package postgresdao

type Organization struct {
	OrgId            int
	OrganizationName string
}

func (dao *DAO) GetOrganizationById(orgId int) (*Organization, error) {
	organization := &Organization{}
	err := dao.db.QueryRow("SELECT org_id, organization_name FROM organizations WHERE org_id = $1", orgId).Scan(&organization.OrgId, &organization.OrganizationName)
	if err != nil {
		return nil, err
	}
	return organization, nil
}

func (dao *DAO) AddOrganization(organization *Organization) error {
	_, err := dao.db.Exec("INSERT INTO organizations (organization_name) VALUES ($1)", organization.OrganizationName)
	return err
}
