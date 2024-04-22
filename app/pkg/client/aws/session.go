package aws_client

import (
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// GetAWSSession returns a singleton AWS session instance
func GetAWSSession() *session.Session {
	var once sync.Once
	var sess *session.Session
	var initErr error
	once.Do(func() {
		// Initialize a session that the SDK uses to load credentials
		// from the shared credentials file ~/.aws/credentials and
		// region from the shared configuration file ~/.aws/config.
		sess, initErr = session.NewSession(&aws.Config{
			Region: aws.String("us-west-2"), // Specify your region
		})
		if initErr != nil {
			fmt.Printf("Failed to create session: %s\n", initErr)
			return
		}
	})
	return sess
}
