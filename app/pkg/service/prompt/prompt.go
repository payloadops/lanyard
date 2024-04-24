package promptservice

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"

	awsclient "plato/app/pkg/client/aws"
	dbdal "plato/app/pkg/dal/postgres"
	"plato/app/pkg/model"
	"plato/app/pkg/util"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

type Service struct {
	s3Client *s3.Client
}

var PROMPT_KEY = "%s/%s/%s/prompt.txt"
var PROMPT_STUB_SIZE = 100

// NewService initializes a new prompt service using AWS S3 client.
func NewService() (*Service, error) {
	return &Service{s3Client: awsclient.GetS3Client()}, nil
}

// GetPrompt retrieves a prompt from S3 based on the provided identifiers.
func (s *Service) GetPrompt(
	ctx context.Context,
	promptId string,
	branch string,
) (*model.GetPromptResponse, error) {
	orgId, orgIdOk := ctx.Value("org_id").(string)
	teamId, teamIdOk := ctx.Value("team_id").(string)
	projectId, projectIdOk := ctx.Value("project_id").(string)

	// Check if all required context values are successfully retrieved
	if !orgIdOk || !teamIdOk || !projectIdOk {
		return nil, fmt.Errorf("failed to parse ids from context")
	}

	key := fmt.Sprintf(PROMPT_KEY, projectId, promptId, branch)
	obj, err := s.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(getBucketForOrgTeamIds(orgId, teamId)),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w.", err)
	}
	defer obj.Body.Close()

	promptBytes, err := io.ReadAll(obj.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read object data: %w", err)
	}

	prompt := string(promptBytes)
	response := &model.GetPromptResponse{
		Prompt:    prompt,
		PromptId:  promptId,
		ProjectId: projectId,
		Branch:    branch,
		Version:   *obj.VersionId,
		Stub:      buildPromptStub(prompt),
	}

	return response, nil
}

// CreatePrompt creates a new prompt in S3.
func (s *Service) CreatePrompt(
	ctx context.Context,
	prompt string,
	branch string,
) (*model.GetPromptResponse, error) {
	orgId, orgIdOk := ctx.Value("org_id").(string)
	teamId, teamIdOk := ctx.Value("team_id").(string)
	projectId, projectIdOk := ctx.Value("project_id").(string)
	promptId := util.GenUUIDString()

	// Check if all required context values are successfully retrieved
	if !orgIdOk || !teamIdOk || !projectIdOk {
		return nil, fmt.Errorf("failed to parse ids from context")
	}

	key := fmt.Sprintf(PROMPT_KEY, projectId, promptId, branch)

	// Attempt to put the prompt into an S3 bucket
	obj, s3Err := s.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(getBucketForOrgTeamIds(orgId, teamId)),
		Key:         aws.String(key),
		Body:        bytes.NewReader([]byte(prompt)),
		ContentType: aws.String("text/plain"),
	})

	if s3Err != nil {
		return nil, fmt.Errorf("error uploading prompt to S3: %w", s3Err)
	}

	// Attempt to add the prompt to the database
	stub := buildPromptStub(prompt)
	_, dbErr := dbdal.AddPrompt(ctx, stub, projectId, promptId, key, *obj.VersionId)
	if dbErr != nil {
		return nil, fmt.Errorf("error recording prompt in database: %w", dbErr)
	}

	response := &model.GetPromptResponse{
		Prompt:    prompt,
		PromptId:  promptId,
		ProjectId: projectId,
		Branch:    branch,
		Version:   *obj.VersionId,
		Stub:      stub,
	}

	return response, nil
}

// DeletePrompt deletes a prompt from S3.
func (s *Service) DeletePrompt(
	ctx context.Context,
	promptId string,
	branch string,
) (*model.DeletePromptResponse, error) {
	projectId, projectIdOk := ctx.Value("project_id").(string)

	// Check if all required context values are successfully retrieved
	if !projectIdOk {
		return nil, fmt.Errorf("failed to parse project id from context")
	}

	err := dbdal.UpdatePromptDeletedStatus(ctx, promptId, true)
	if err != nil {
		return nil, err
	}

	response := &model.DeletePromptResponse{
		PromptId:  promptId,
		ProjectId: projectId,
	}
	return response, err
}

// ListPrompts lists all prompts for a given project.
func (s *Service) ListPrompts(
	ctx context.Context,
	orgId string,
	teamId string,
	projectId string,
) error {
	var err error
	if err != nil {
		return fmt.Errorf("failed to list prompts: %w", err)
	}
	return nil
}

func buildPromptStub(prompt string) string {
	if len(prompt) > PROMPT_STUB_SIZE {
		return prompt[:PROMPT_STUB_SIZE] + "..."
	}
	return prompt
}

func getBucketForOrgTeamIds(orgId string, teamId string) string {
	hasher := md5.New()

	_, err := hasher.Write([]byte(fmt.Sprintf("%s/%s", orgId, teamId)))
	if err != nil {
		panic(err)
	}

	hashBytes := hasher.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)

	return hashString
}
