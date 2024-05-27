import * as codebuild from 'aws-cdk-lib/aws-codebuild';

export const BUILDSPEC = codebuild.BuildSpec.fromObject({
  version: '0.2',
  phases: {
    pre_build: {
      commands: [
        'echo Logging in to Amazon ECR...',
        '$(aws ecr get-login --no-include-email --region $AWS_DEFAULT_REGION)'
      ]
    },
    build: {
      commands: [
        'echo Building the Docker image from the app directory...',
        'cd app',  // Navigate to the app directory where the Dockerfile is located
        'docker build -t $REPOSITORY_URI:$CODEBUILD_RESOLVED_SOURCE_VERSION .',
        'docker tag $REPOSITORY_URI:$CODEBUILD_RESOLVED_SOURCE_VERSION $REPOSITORY_URI:latest'
      ]
    },
    post_build: {
      commands: [
        'echo Pushing the Docker image...',
        'docker push $REPOSITORY_URI:$CODEBUILD_RESOLVED_SOURCE_VERSION',
        'docker push $REPOSITORY_URI:latest',
        'echo Write the image definitions file...',
        'printf \'[{"name":"my-container","imageUri":"%s"}]\' $REPOSITORY_URI:$CODEBUILD_RESOLVED_SOURCE_VERSION > imagedefinitions.json'
      ]
    }
  },
  artifacts: {
    files: ['imagedefinitions.json']
  },
  cache: {
    paths: [
      '/root/.m2/**/*',
      '/root/.gradle/**/*'
    ]
  }
});
