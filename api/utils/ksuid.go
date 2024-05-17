package utils

import (
	"github.com/segmentio/ksuid"
)

// GenerateKSUID generates a new KSUID, which is unique and sortable.
func GenerateKSUID() (string, error) {
	id := ksuid.New()
	return id.String(), nil
}
