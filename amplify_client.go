package amplifyx

import (
	"context"

	awsamplify "github.com/aws/aws-sdk-go-v2/service/amplify"
)

type AmplifyClient interface {
	StartJob(ctx context.Context, params *awsamplify.StartJobInput, optFns ...func(*awsamplify.Options)) (*awsamplify.StartJobOutput, error)
	GetJob(ctx context.Context, params *awsamplify.GetJobInput, optFns ...func(*awsamplify.Options)) (*awsamplify.GetJobOutput, error)
}
