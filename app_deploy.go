package amplifyx

import (
	"context"
	"slices"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsamplify "github.com/aws/aws-sdk-go-v2/service/amplify"
	awsamplifytypes "github.com/aws/aws-sdk-go-v2/service/amplify/types"
	"github.com/samber/oops"
)

// DeployArgs represents the arguments of the deploy command.
type DeployArgs struct {
	AppID               string        `required:""`
	Branch              string        `default:"main"`
	ObservationInterval time.Duration `default:"5s"`
}

// Validate returns an error if args is invalid.
func (args DeployArgs) Validate() error {
	errs := []error{}

	if len(args.AppID) == 0 {
		errs = append(errs, oops.New("invalid app id: must not be empty"))
	}
	if len(args.Branch) == 0 {
		errs = append(errs, oops.New("invalid branch: must not be empty"))
	}
	if args.ObservationInterval < time.Second {
		errs = append(errs, oops.New("invalid observation interval: must be at least 1s"))
	}

	return oops.Join(errs...)
}

// Deploy starts a release job for the branch of the Amplify app and waits until the job finishes.
// It returns an error if the job fails or is cancelled.
func (app App) Deploy(ctx context.Context, args DeployArgs) error {
	var jobSummary *awsamplifytypes.JobSummary
	{
		out, err := app.amplifyClient.StartJob(ctx, &awsamplify.StartJobInput{
			AppId:      aws.String(args.AppID),
			BranchName: aws.String(args.Branch),
			JobType:    awsamplifytypes.JobTypeRelease,
		})
		if err != nil {
			return oops.Wrapf(err, "failed to start job")
		}

		jobSummary = out.JobSummary
	}

	for slices.Contains([]awsamplifytypes.JobStatus{
		awsamplifytypes.JobStatusCreated,
		awsamplifytypes.JobStatusPending,
		awsamplifytypes.JobStatusProvisioning,
		awsamplifytypes.JobStatusRunning,
		awsamplifytypes.JobStatusCancelling,
	}, jobSummary.Status) {
		select {
		case <-ctx.Done():
			return oops.Wrap(context.Cause(ctx))
		case <-time.After(args.ObservationInterval):
		}

		{
			out, err := app.amplifyClient.GetJob(ctx, &awsamplify.GetJobInput{
				AppId:      aws.String(args.AppID),
				BranchName: aws.String(args.Branch),
				JobId:      jobSummary.JobId,
			})
			if err != nil {
				return oops.Wrapf(err, "failed to get job")
			}

			jobSummary = out.Job.Summary
		}
	}

	switch jobStatus := jobSummary.Status; jobStatus {
	case awsamplifytypes.JobStatusSucceed:
		return nil
	case awsamplifytypes.JobStatusFailed:
		return oops.New("failed to deploy: job failed")
	case awsamplifytypes.JobStatusCancelled:
		return oops.New("failed to deploy: job cancelled")
	default:
		return oops.Errorf("unsupported job status: %s", jobStatus)
	}
}
