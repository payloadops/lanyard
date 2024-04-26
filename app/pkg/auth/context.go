package auth

import "github.com/dgrijalva/jwt-go"

type ProjectId struct{}
type UserId struct{}
type OrgId struct{}
type Claims struct {
	UserId string `json:"user_id"`
	OrgId  string `json:"org_id"` // This could also be `sub` if the user Id is stored in the subject claim
	jwt.StandardClaims
}
