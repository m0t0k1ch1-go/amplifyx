package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/alecthomas/kong"
	"github.com/samber/oops"

	"github.com/m0t0k1ch1-go/amplifyx/v2"
)

var (
	cmd  string
	args amplifyx.Args
)

func init() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, nil)))
}

func main() {
	os.Exit(run(context.Background()))
}

func run(ctx context.Context) int {
	if err := parse(); err != nil {
		return fail(ctx, err)
	}

	if args.Timeout > 0 {
		var cancel context.CancelFunc

		ctx, cancel = context.WithTimeout(ctx, args.Timeout)
		defer cancel()
	}

	var err error

	if panicErr := oops.Recover(func() {
		err = command(ctx)
	}); panicErr != nil {
		err = panicErr
	}

	if err != nil {
		return fail(ctx, err)
	}

	return 0
}

func parse() error {
	k, err := kong.New(&args)
	if err != nil {
		return oops.Wrapf(err, "failed to initialize args parser")
	}

	kctx, err := k.Parse(os.Args[1:])
	if err != nil {
		return oops.Wrapf(err, "failed to parse args")
	}

	cmd = kctx.Command()

	return nil
}

func command(ctx context.Context) error {
	app, err := amplifyx.NewApp(ctx)
	if err != nil {
		return oops.Wrapf(err, "failed to initialize app")
	}

	switch cmd {
	case "deploy":
		return app.Deploy(ctx, args.Deploy)
	case "version":
		return app.Version(ctx, args.Version)
	default:
		return oops.Errorf("unexpected command: %s", cmd)
	}
}

func fail(ctx context.Context, err error) int {
	slog.ErrorContext(ctx, err.Error())

	return 1
}
