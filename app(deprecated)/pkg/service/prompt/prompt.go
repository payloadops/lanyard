package promptservice

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"plato/app/pkg/auth"
	awsclient "plato/app/pkg/client/aws"
	dbdal "plato/app/pkg/dal/postgres"
	promptservicemodel "plato/app/pkg/model/prompt/service"
	"plato/app/pkg/util"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

var _ PromptService = (*Service)(nil)

type Service struct {
	s3Client *s3.Client
}

var PROMPT_KEY = "%s/%s/prompt.txt"
var PROMPT_STUB_SIZE = 100

// NewService initializes a new prompt service using AWS S3 client.
func NewService() (PromptService, error) {
	return &Service{s3Client: awsclient.GetS3Client()}, nil
}

// GetPrompt retrieves a prompt from S3 based on the provided identifiers.
func (s *Service) GetPrompt(
	ctx context.Context,
	projectId string,
	promptId string,
	branch string,
	version string,
) (*promptservicemodel.GetPromptResponse, error) {
	orgId, orgIdOk := ctx.Value(auth.OrgContext{}).(string)

	// Check if all required context values are successfully retrieved
	if !orgIdOk {
		return nil, fmt.Errorf("failed to parse ids from context")
	}

	promptRecord, dbErr := dbdal.GetPromptById(ctx, projectId, promptId)

	if dbErr != nil {
		return nil, fmt.Errorf("failed to read object metadata: %w", dbErr)
	}
	if promptRecord.Deleted {
		return nil, fmt.Errorf("prompt cannot be retrieved as its marked as deleted")
	}
	if len(branch) == 0 {
		branch = promptRecord.DefaultBranch
	}

	key := fmt.Sprintf(PROMPT_KEY, promptId, branch)
	getObjectInput := &s3.GetObjectInput{
		Bucket: aws.String(util.GetBucketString(orgId, projectId)),
		Key:    aws.String(key),
	}

	if len(version) > 0 {
		getObjectInput.VersionId = aws.String(version)
	}

	obj, err := s.s3Client.GetObject(ctx, getObjectInput)
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}
	defer obj.Body.Close()

	promptBytes, err := io.ReadAll(obj.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read object data: %w", err)
	}

	prompt := string(promptBytes)
	response := &promptservicemodel.GetPromptResponse{
		Prompt:     prompt,
		Name:       promptRecord.Name,
		PromptId:   promptId,
		ProjectId:  projectId,
		Branch:     branch,
		Version:    *obj.VersionId,
		ModifiedAt: promptRecord.ModifiedAt,
		Stub:       buildPromptStub(prompt),
	}

	return response, nil
}

// CreatePrompt creates a new prompt in S3.
func (s *Service) CreatePrompt(
	ctx context.Context,
	projectId string,
	createPromptRequest promptservicemodel.CreatePromptRequest,
) (*promptservicemodel.GetPromptResponse, error) {
	orgId, orgIdOk := ctx.Value(auth.OrgContext{}).(string)
	promptId := util.GenIDString()

	// Check if all required context values are successfully retrieved
	if !orgIdOk {
		return nil, fmt.Errorf("failed to parse ids from context")
	}

	if len(createPromptRequest.Branch) == 0 {
		return nil, fmt.Errorf("no branch provided")
	}

	key := fmt.Sprintf(PROMPT_KEY, promptId, createPromptRequest.Branch)

	// Attempt to put the prompt into an S3 bucket
	obj, s3Err := s.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(util.GetBucketString(orgId, projectId)),
		Key:         aws.String(key),
		Body:        bytes.NewReader([]byte(createPromptRequest.Prompt)),
		ContentType: aws.String("text/plain"),
	})

	if s3Err != nil {
		return nil, fmt.Errorf("error uploading prompt to S3: %w", s3Err)
	}

	// Attempt to add the prompt to the database
	stub := buildPromptStub(createPromptRequest.Prompt)
	promptRecord, dbErr := dbdal.AddPrompt(ctx, createPromptRequest.Name, stub, projectId, promptId, fmt.Sprintf("%s/%s", projectId, promptId), *obj.VersionId, createPromptRequest.Branch)
	if dbErr != nil {
		return nil, fmt.Errorf("error recording prompt in database: %w", dbErr)
	}

	response := &promptservicemodel.GetPromptResponse{
		Prompt:     createPromptRequest.Prompt,
		Name:       promptRecord.Name,
		PromptId:   promptId,
		ProjectId:  projectId,
		Branch:     createPromptRequest.Branch,
		Version:    *obj.VersionId,
		ModifiedAt: promptRecord.ModifiedAt,
		Stub:       stub,
	}

	return response, nil
}

func (s *Service) UpdatePrompt(
	ctx context.Context,
	projectId string,
	promptId string,
	updatePromptRequest promptservicemodel.UpdatePromptRequest,
) (*promptservicemodel.GetPromptResponse, error) {
	orgId, orgIdOk := ctx.Value(auth.OrgContext{}).(string)

	// Check if all required context values are successfully retrieved
	if !orgIdOk {
		return nil, fmt.Errorf("failed to parse ids from context")
	}

	key := fmt.Sprintf(PROMPT_KEY, promptId, updatePromptRequest.Branch)

	// Attempt to put the prompt into an S3 bucket
	obj, s3Err := s.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(util.GetBucketString(orgId, projectId)),
		Key:         aws.String(key),
		Body:        bytes.NewReader([]byte(updatePromptRequest.Prompt)),
		ContentType: aws.String("text/plain"),
	})

	if s3Err != nil {
		return nil, fmt.Errorf("error uploading prompt to S3: %w", s3Err)
	}

	// Attempt to add the prompt to the database
	stub := buildPromptStub(updatePromptRequest.Prompt)
	modifiedAt, dbErr := dbdal.UpdatePrompt(ctx, updatePromptRequest.Name, projectId, promptId, stub, *obj.VersionId)
	if dbErr != nil {
		return nil, fmt.Errorf("error recording prompt in database: %w", dbErr)
	}

	response := &promptservicemodel.GetPromptResponse{
		Prompt:     updatePromptRequest.Prompt,
		Name:       updatePromptRequest.Name,
		PromptId:   promptId,
		ProjectId:  projectId,
		Branch:     updatePromptRequest.Branch,
		Version:    *obj.VersionId,
		Stub:       stub,
		ModifiedAt: modifiedAt,
	}

	return response, nil
}

// Soft deletes prompt record in DB only.
func (s *Service) DeletePrompt(
	ctx context.Context,
	projectId string,
	promptId string,
) (*promptservicemodel.DeletePromptResponse, error) {
	deletedAt, err := dbdal.UpdatePromptDeletedStatus(ctx, projectId, promptId, true)
	if err != nil {
		return nil, err
	}

	response := &promptservicemodel.DeletePromptResponse{
		PromptId:  promptId,
		ProjectId: projectId,
		DeletedAt: deletedAt,
	}
	return response, err
}

// ListPrompts lists all prompts for a given project.
func (s *Service) ListPrompts(
	ctx context.Context,
	projectId string,
) (*promptservicemodel.ListPromptsResponse, error) {
	var err error
	prompts, err := dbdal.ListPromptsByProjectId(ctx, projectId)
	if err != nil {
		return nil, fmt.Errorf("failed to list prompts: %w", err)
	}

	response := &promptservicemodel.ListPromptsResponse{
		Prompts:   &prompts,
		ProjectId: projectId,
	}

	return response, nil
}

func (s *Service) ListVersions(
	ctx context.Context,
	projectId string,
	promptId string,
) (*promptservicemodel.ListVersionsResponse, error) {
	var err error
	_, err = dbdal.ListPromptsByProjectId(ctx, projectId)
	if err != nil {
		return nil, fmt.Errorf("failed to list prompts: %w", err)
	}

	response := &promptservicemodel.ListVersionsResponse{}

	return response, nil
}

func (s *Service) UpdateActiveVersion(
	ctx context.Context,
	projectId string,
	promptId string,
	updateActiveVersionRequest *promptservicemodel.UpdateActiveVersionRequest,
) (*promptservicemodel.GetPromptResponse, error) {
	orgId, orgIdOk := ctx.Value(auth.OrgContext{}).(string)

	if !orgIdOk {
		return nil, fmt.Errorf("failed to parse ids from context")
	}

	key := fmt.Sprintf(PROMPT_KEY, promptId, updateActiveVersionRequest.Branch)
	getObj, s3GetErr := s.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket:    aws.String(util.GetBucketString(orgId, projectId)),
		Key:       aws.String(key),
		VersionId: &updateActiveVersionRequest.Version,
	})

	if s3GetErr != nil {
		return nil, fmt.Errorf("failed to get object: %w", s3GetErr)
	}
	defer getObj.Body.Close()

	promptBytes, err := io.ReadAll(getObj.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read object data: %w", err)
	}

	putObj, s3PutErr := s.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(util.GetBucketString(orgId, projectId)),
		Key:         aws.String(key),
		Body:        bytes.NewReader(promptBytes),
		ContentType: aws.String("text/plain"),
	})

	if s3PutErr != nil {
		return nil, fmt.Errorf("error uploading prompt to S3: %w", s3PutErr)
	}

	prompt := string(promptBytes)

	stub := buildPromptStub(prompt)
	modifiedAt, dbErr := dbdal.UpdatePromptActiveVersion(ctx, projectId, promptId, stub, *putObj.VersionId)
	if dbErr != nil {
		return nil, fmt.Errorf("error recording prompt in database: %w", dbErr)
	}

	response := &promptservicemodel.GetPromptResponse{
		Prompt:     prompt,
		PromptId:   promptId,
		ProjectId:  projectId,
		Branch:     updateActiveVersionRequest.Branch,
		Version:    *putObj.VersionId,
		Stub:       stub,
		ModifiedAt: modifiedAt,
	}

	return response, nil
}

func (s *Service) CreateBranch(
	ctx context.Context,
	projectId string,
	promptId string,
	createBranchRequest promptservicemodel.CreateBranchRequest,
) (*promptservicemodel.CreateBranchResponse, error) {
	var err error
	_, err = dbdal.ListPromptsByProjectId(ctx, projectId)
	if err != nil {
		return nil, fmt.Errorf("failed to list prompts: %w", err)
	}

	response := &promptservicemodel.CreateBranchResponse{}

	return response, nil
}

func (s *Service) DeleteBranch(
	ctx context.Context,
	projectId string,
	promptId string,
	branch string,
) (*promptservicemodel.DeleteBranchResponse, error) {
	var err error
	_, err = dbdal.ListPromptsByProjectId(ctx, projectId)
	if err != nil {
		return nil, fmt.Errorf("failed to list prompts: %w", err)
	}

	response := &promptservicemodel.DeleteBranchResponse{}

	return response, nil
}

func (s *Service) ListBranches(
	ctx context.Context,
	projectId string,
	promptId string,
) (*promptservicemodel.ListBranchesResponse, error) {
	var err error
	_, err = dbdal.ListPromptsByProjectId(ctx, projectId)
	if err != nil {
		return nil, fmt.Errorf("failed to list prompts: %w", err)
	}

	response := &promptservicemodel.ListBranchesResponse{}

	return response, nil
}

func buildPromptStub(prompt string) string {
	if len(prompt) > PROMPT_STUB_SIZE {
		return prompt[:PROMPT_STUB_SIZE] + "..."
	}
	return prompt
}
