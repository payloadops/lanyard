package util

import "github.com/google/uuid"

func GenUUIDString() string {
	uuid, _ := uuid.NewRandom()
	return uuid.String()
}
