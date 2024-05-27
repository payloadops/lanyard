package util

import "github.com/rs/xid"

// GenIDString generates and returns a globally unique identifier as a string.
func GenIDString() string {
	id := xid.New()
	return id.String()
}
