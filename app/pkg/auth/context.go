package auth

// ProjectContext represents the context specific to a project,
// containing project-specific identifiers.
type ProjectContext struct {
	Id string // Id is the unique identifier of the project.
}

// UserContext represents the context specific to a user,
// containing user-specific identifiers.
type UserContext struct {
	Id string // Id is the unique identifier of the user.
}

// OrgContext represents the context specific to an organization,
// containing organization-specific identifiers.
type OrgContext struct {
	Id string // Id is the unique identifier of the organization.
}
