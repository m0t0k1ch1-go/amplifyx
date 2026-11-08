package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/alecthomas/kong"
	"github.com/samber/oops"

	"github.com/m0t0k1ch1-go/amplifyx/v2"
)

func init() {
	slog.SetDefault(amplifyx.NewLogger(os.Stdout, os.Stderr))
}

func main() {
	os.Exit(run(context.Background()))
}

func run(ctx context.Context) int {
	var (
		command string
		args    amplifyx.Args
	)
	{
		k, err := kong.New(&args)
		if err != nil {
			return fail(ctx, oops.Wrapf(err, "failed to initialize args parser"))
		}

		kctx, err := k.Parse(os.Args[1:])
		if err != nil {
			return fail(ctx, oops.Wrapf(err, "failed to parse args"))
		}

		command = kctx.Command()
	}

	if args.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, args.Timeout)
		defer cancel()
	}

	app, err := amplifyx.InitializeApp(ctx)
	if err != nil {
		return fail(ctx, oops.Wrapf(err, "failed to initialize app"))
	}

	{
		var err error

		if panicErr := oops.Recover(func() {
			switch command {
			case "deploy":
				err = app.Deploy(ctx, args.Deploy)
			case "version":
				err = app.Version(ctx, args.Version)
			default:
				err = oops.Errorf("unsupported command: %s", command)
			}
		}); panicErr != nil {
			err = panicErr
		}
		if err != nil {
			return fail(ctx, err)
		}
	}

	return 0
}

func fail(ctx context.Context, err error) int {
	slog.ErrorContext(ctx, "failed to run", slog.Any("error", err))

	return 1
}
