package util

import (
	"crypto/md5"
	"fmt"
	"io"
)

func GetBucketString(orgId string, projectId string) string {
	hasher := md5.New()

	str := fmt.Sprintf("%s-%s", orgId, projectId)
	io.WriteString(hasher, str)
	if _, err := io.WriteString(hasher, str); err != nil {
		panic("Failed to build bucket hash")
	}

	md5hash := hasher.Sum(nil)

	return string(md5hash)
}
