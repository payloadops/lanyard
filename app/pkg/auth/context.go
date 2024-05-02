package auth

import "github.com/dgrijalva/jwt-go"

type ProjectContext struct {
	Id string
}

type UserContext struct {
	Id string
}

type OrgContext struct {
	Id string
}

type Claims struct {
	jwt.StandardClaims
	UserId string `json:"user_id"`
	OrgId  string `json:"org_id"`
}
