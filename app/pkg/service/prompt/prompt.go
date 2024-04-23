package promptservice

import (
	"context"
	"fmt"
	"io"
	awsclient "plato/app/pkg/client/aws"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

type Service struct {
	s3Client *s3.Client
}

func NewService() (*Service, error) {
	return &Service{s3Client: awsclient.GetS3Client()}, nil
}

type GetPromptRequest struct {
	Bucket string
	Key    string
}

type GetPromptResponse struct {
	Prompt string `json:"prompt"`
}

// GetPrompt retrieves a prompt from S3 based on the specified bucket and key.
func (s *Service) GetPrompt(ctx context.Context, req GetPromptRequest) (*GetPromptResponse, error) {
	obj, err := s.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(req.Bucket),
		Key:    aws.String(req.Key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}
	defer obj.Body.Close()

	promptBytes, err := io.ReadAll(obj.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read object data: %w", err)
	}

	return &GetPromptResponse{Prompt: string(promptBytes)}, nil
}
