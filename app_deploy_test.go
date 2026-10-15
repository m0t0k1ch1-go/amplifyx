package amplifyx_test

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/alecthomas/kong"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsamplify "github.com/aws/aws-sdk-go-v2/service/amplify"
	awsamplifytypes "github.com/aws/aws-sdk-go-v2/service/amplify/types"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/require"

	"github.com/m0t0k1ch1-go/amplifyx/v2"
	"github.com/m0t0k1ch1-go/amplifyx/v2/internal/gen/moq"
	"github.com/m0t0k1ch1-go/amplifyx/v2/internal/testutil"
)

func TestDeployArgs_Validate(t *testing.T) {
	t.Run("failure", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []string
			want []string
		}{
			{
				"with all arguments invalid",
				[]string{"deploy", "--branch=", "--observation-interval=999ms"},
				[]string{
					"invalid app id: must not be empty",
					"invalid branch: must not be empty",
					"invalid observation interval: must be at least 1s",
				},
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var args amplifyx.Args

				k, err := kong.New(&args)
				require.NoError(t, err)

				_, err = k.Parse(tc.in)
				for _, want := range tc.want {
					require.ErrorContains(t, err, want)
				}
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []string
			want amplifyx.DeployArgs
		}{
			{
				"with the defaults",
				[]string{"deploy", "--app-id=app"},
				amplifyx.DeployArgs{
					AppID:               "app",
					Branch:              "main",
					ObservationInterval: 5 * time.Second,
				},
			},
			{
				"with all arguments valid",
				[]string{"deploy", "--app-id=app", "--branch=dev", "--observation-interval=1s"},
				amplifyx.DeployArgs{
					AppID:               "app",
					Branch:              "dev",
					ObservationInterval: time.Second,
				},
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var args amplifyx.Args

				k, err := kong.New(&args)
				require.NoError(t, err)

				_, err = k.Parse(tc.in)
				require.NoError(t, err)
				testutil.RequireEqual(t, tc.want, args.Deploy)
			})
		}
	})
}

func TestApp_Deploy(t *testing.T) {
	newAmplifyStartJobFunc := func(appID string, branch string, jobID string) func(ctx context.Context, params *awsamplify.StartJobInput, optFns ...func(*awsamplify.Options)) (*awsamplify.StartJobOutput, error) {
		return func(ctx context.Context, params *awsamplify.StartJobInput, optFns ...func(*awsamplify.Options)) (*awsamplify.StartJobOutput, error) {
			if err := testutil.Equal(&awsamplify.StartJobInput{
				AppId:      aws.String(appID),
				BranchName: aws.String(branch),
				JobType:    awsamplifytypes.JobTypeRelease,
			}, params, cmpopts.IgnoreUnexported(awsamplify.StartJobInput{})); err != nil {
				return nil, fmt.Errorf("invalid params: %w", err)
			}

			return &awsamplify.StartJobOutput{
				JobSummary: &awsamplifytypes.JobSummary{
					JobId:  aws.String(jobID),
					Status: awsamplifytypes.JobStatusCreated,
				},
			}, nil
		}
	}
	newAmplifyGetJobFunc := func(appID string, branch string, jobID string, jobStatuses []awsamplifytypes.JobStatus) func(ctx context.Context, params *awsamplify.GetJobInput, optFns ...func(*awsamplify.Options)) (*awsamplify.GetJobOutput, error) {
		callCnt := 0

		return func(ctx context.Context, params *awsamplify.GetJobInput, optFns ...func(*awsamplify.Options)) (*awsamplify.GetJobOutput, error) {
			callCnt++
			if callCnt > len(jobStatuses) {
				return nil, errors.New("unsupported call")
			}

			if err := testutil.Equal(&awsamplify.GetJobInput{
				AppId:      aws.String(appID),
				BranchName: aws.String(branch),
				JobId:      aws.String(jobID),
			}, params, cmpopts.IgnoreUnexported(awsamplify.GetJobInput{})); err != nil {
				return nil, fmt.Errorf("invalid params: %w", err)
			}

			return &awsamplify.GetJobOutput{
				Job: &awsamplifytypes.Job{
					Summary: &awsamplifytypes.JobSummary{
						JobId:  aws.String(jobID),
						Status: jobStatuses[callCnt-1],
					},
				},
			}, nil
		}
	}

	t.Run("failure: failed to start job", func(t *testing.T) {
		ctx := t.Context()

		args := amplifyx.DeployArgs{
			AppID:  "app",
			Branch: "main",
		}

		amplifyClient := &moq.AmplifyClientMock{
			StartJobFunc: func(ctx context.Context, params *awsamplify.StartJobInput, optFns ...func(*awsamplify.Options)) (*awsamplify.StartJobOutput, error) {
				if err := testutil.Equal(&awsamplify.StartJobInput{
					AppId:      aws.String(args.AppID),
					BranchName: aws.String(args.Branch),
					JobType:    awsamplifytypes.JobTypeRelease,
				}, params, cmpopts.IgnoreUnexported(awsamplify.StartJobInput{})); err != nil {
					return nil, fmt.Errorf("invalid params: %w", err)
				}

				return nil, errors.New("something went wrong")
			},
		}

		app := amplifyx.NewApp(io.Discard, amplifyClient)

		err := app.Deploy(ctx, args)
		require.ErrorContains(t, err, "failed to start job: something went wrong")

		require.Len(t, amplifyClient.StartJobCalls(), 1)
	})

	t.Run("failure: failed to get job", func(t *testing.T) {
		ctx := t.Context()

		args := amplifyx.DeployArgs{
			AppID:               "app",
			Branch:              "main",
			ObservationInterval: time.Millisecond,
		}

		jobID := rand.Text()

		amplifyClient := &moq.AmplifyClientMock{
			StartJobFunc: newAmplifyStartJobFunc(args.AppID, args.Branch, jobID),
			GetJobFunc: func(ctx context.Context, params *awsamplify.GetJobInput, optFns ...func(*awsamplify.Options)) (*awsamplify.GetJobOutput, error) {
				if err := testutil.Equal(&awsamplify.GetJobInput{
					AppId:      aws.String(args.AppID),
					BranchName: aws.String(args.Branch),
					JobId:      aws.String(jobID),
				}, params, cmpopts.IgnoreUnexported(awsamplify.GetJobInput{})); err != nil {
					return nil, fmt.Errorf("invalid params: %w", err)
				}

				return nil, errors.New("something went wrong")
			},
		}

		app := amplifyx.NewApp(io.Discard, amplifyClient)

		err := app.Deploy(ctx, args)
		require.ErrorContains(t, err, "failed to get job: something went wrong")

		require.Len(t, amplifyClient.StartJobCalls(), 1)
		require.Len(t, amplifyClient.GetJobCalls(), 1)
	})

	t.Run("failure: failed to deploy: job failed", func(t *testing.T) {
		ctx := t.Context()

		args := amplifyx.DeployArgs{
			AppID:               "app",
			Branch:              "main",
			ObservationInterval: time.Millisecond,
		}

		jobID := rand.Text()

		amplifyClient := &moq.AmplifyClientMock{}
		{
			amplifyClient.StartJobFunc = newAmplifyStartJobFunc(args.AppID, args.Branch, jobID)
			amplifyClient.GetJobFunc = newAmplifyGetJobFunc(args.AppID, args.Branch, jobID, []awsamplifytypes.JobStatus{
				awsamplifytypes.JobStatusPending,
				awsamplifytypes.JobStatusFailed,
			})
		}

		app := amplifyx.NewApp(io.Discard, amplifyClient)

		err := app.Deploy(ctx, args)
		require.ErrorContains(t, err, "failed to deploy: job failed")

		require.Len(t, amplifyClient.StartJobCalls(), 1)
		require.Len(t, amplifyClient.GetJobCalls(), 2)
	})

	t.Run("failure: failed to deploy: job canceled", func(t *testing.T) {
		ctx := t.Context()

		args := amplifyx.DeployArgs{
			AppID:               "app",
			Branch:              "main",
			ObservationInterval: time.Millisecond,
		}

		jobID := rand.Text()

		amplifyClient := &moq.AmplifyClientMock{}
		{
			amplifyClient.StartJobFunc = newAmplifyStartJobFunc(args.AppID, args.Branch, jobID)
			amplifyClient.GetJobFunc = newAmplifyGetJobFunc(args.AppID, args.Branch, jobID, []awsamplifytypes.JobStatus{
				awsamplifytypes.JobStatusPending,
				awsamplifytypes.JobStatusCancelled,
			})
		}

		app := amplifyx.NewApp(io.Discard, amplifyClient)

		err := app.Deploy(ctx, args)
		require.ErrorContains(t, err, "failed to deploy: job canceled")

		require.Len(t, amplifyClient.StartJobCalls(), 1)
		require.Len(t, amplifyClient.GetJobCalls(), 2)
	})

	t.Run("failure: unsupported job status", func(t *testing.T) {
		ctx := t.Context()

		args := amplifyx.DeployArgs{
			AppID:               "app",
			Branch:              "main",
			ObservationInterval: time.Millisecond,
		}

		jobID := rand.Text()
		jobStatus := awsamplifytypes.JobStatus("UNKNOWN")

		amplifyClient := &moq.AmplifyClientMock{
			StartJobFunc: newAmplifyStartJobFunc(args.AppID, args.Branch, jobID),
			GetJobFunc: newAmplifyGetJobFunc(args.AppID, args.Branch, jobID, []awsamplifytypes.JobStatus{
				jobStatus,
			}),
		}

		app := amplifyx.NewApp(io.Discard, amplifyClient)

		err := app.Deploy(ctx, args)
		require.ErrorContains(t, err, "unsupported job status: "+string(jobStatus))

		require.Len(t, amplifyClient.StartJobCalls(), 1)
		require.Len(t, amplifyClient.GetJobCalls(), 1)
	})

	t.Run("failure: context canceled", func(t *testing.T) {
		ctx := t.Context()
		cctx, cancel := context.WithCancel(ctx)

		args := amplifyx.DeployArgs{
			AppID:               "app",
			Branch:              "main",
			ObservationInterval: time.Minute,
		}

		amplifyClient := &moq.AmplifyClientMock{
			StartJobFunc: newAmplifyStartJobFunc(args.AppID, args.Branch, rand.Text()),
		}

		app := amplifyx.NewApp(io.Discard, amplifyClient)

		cancel()

		err := app.Deploy(cctx, args)
		require.ErrorIs(t, err, context.Canceled)

		require.Len(t, amplifyClient.StartJobCalls(), 1)
	})

	t.Run("success", func(t *testing.T) {
		ctx := t.Context()

		args := amplifyx.DeployArgs{
			AppID:               "app",
			Branch:              "main",
			ObservationInterval: time.Millisecond,
		}

		jobID := rand.Text()

		amplifyClient := &moq.AmplifyClientMock{}
		{
			amplifyClient.StartJobFunc = newAmplifyStartJobFunc(args.AppID, args.Branch, jobID)
			amplifyClient.GetJobFunc = newAmplifyGetJobFunc(args.AppID, args.Branch, jobID, []awsamplifytypes.JobStatus{
				awsamplifytypes.JobStatusPending,
				awsamplifytypes.JobStatusProvisioning,
				awsamplifytypes.JobStatusRunning,
				awsamplifytypes.JobStatusSucceed,
			})
		}

		app := amplifyx.NewApp(io.Discard, amplifyClient)

		err := app.Deploy(ctx, args)
		require.NoError(t, err)

		require.Len(t, amplifyClient.StartJobCalls(), 1)
		require.Len(t, amplifyClient.GetJobCalls(), 4)
	})
}
