package util

import "github.com/rs/xid"

func GenIDString() string {
	id := xid.New()
	return id.String()
}
