package util

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

// GetBucketString returns an MD5 hash as a hex string of the concatenated orgId and projectId.
func GetBucketString(orgId string, projectId string) string {
	hasher := md5.New()
	str := fmt.Sprintf("%s-%s", orgId, projectId)
	hasher.Write([]byte(str))
	hashBytes := hasher.Sum(nil)
	return hex.EncodeToString(hashBytes)
}
