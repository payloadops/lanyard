package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetBucketString(t *testing.T) {
	// Test case 1: orgId and projectId are non-empty strings
	orgId := "org1"
	projectId := "project1"
	expectedHash := "02d47f644d36e93a98b9ace246415736"
	resultHash := GetBucketString(orgId, projectId)
	assert.Equal(t, expectedHash, resultHash, "Test case 1 failed")

	// Test case 2: orgId and projectId are empty strings
	orgId = ""
	projectId = ""
	expectedHash = "336d5ebc5436534e61d16e63ddfca327" // MD5 hash of empty string
	resultHash = GetBucketString(orgId, projectId)
	assert.Equal(t, expectedHash, resultHash, "Test case 2 failed")

	// Test case 3: orgId is empty, projectId is non-empty
	orgId = ""
	projectId = "project2"
	expectedHash = "624f0d86fccb8f8054b918e807a64cff"
	resultHash = GetBucketString(orgId, projectId)
	assert.Equal(t, expectedHash, resultHash, "Test case 3 failed")

	// Test case 4: orgId is non-empty, projectId is empty
	orgId = "org2"
	projectId = ""
	expectedHash = "ae39243e5df0ea64129589e30f1a46a9"
	resultHash = GetBucketString(orgId, projectId)
	assert.Equal(t, expectedHash, resultHash, "Test case 4 failed")
}
