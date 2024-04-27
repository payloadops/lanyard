package util

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

func GetBucketString(orgId string, projectId string) string {
	hasher := md5.New()

	str := fmt.Sprintf("%s-%s", orgId, projectId)
	hasher.Write([]byte(str))
	hashBytes := hasher.Sum(nil)
	return hex.EncodeToString(hashBytes)
}
