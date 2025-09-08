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

type DeployArgs struct {
	AppID               string        `required:""`
	Branch              string        `default:"main"`
	ObservationInterval time.Duration `default:"5s"`
}

func (app *App) Deploy(ctx context.Context, args DeployArgs) error {
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
			return context.Cause(ctx)
		default:
		}

		time.Sleep(args.ObservationInterval)

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
		return oops.New("job failed")
	case awsamplifytypes.JobStatusCancelled:
		return oops.New("job cancelled")
	default:
		return oops.Errorf("unexpected job status: %s", jobStatus)
	}
}
