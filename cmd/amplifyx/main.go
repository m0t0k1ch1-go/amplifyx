package main

import (
	"context"
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/samber/oops"

	"github.com/m0t0k1ch1-go/amplifyx/v2"
)

func main() {
	os.Exit(run(context.Background()))
}

func run(ctx context.Context) int {
	var (
		cmd  string
		args amplifyx.Args
	)
	{
		k, err := kong.New(&args)
		if err != nil {
			return fail(oops.Wrapf(err, "failed to initialize args parser"))
		}

		kctx, err := k.Parse(os.Args[1:])
		if err != nil {
			return fail(oops.Wrapf(err, "failed to parse args"))
		}

		cmd = kctx.Command()
	}

	app, err := amplifyx.NewApp(ctx)
	if err != nil {
		return fail(oops.Wrapf(err, "failed to initialize app"))
	}

	if args.Timeout > 0 {
		var cancel context.CancelFunc

		ctx, cancel = context.WithTimeout(ctx, args.Timeout)
		defer cancel()
	}

	if err := command(ctx, app, cmd, args); err != nil {
		return fail(oops.Wrapf(err, "failed to execute command"))
	}

	return 0
}

func command(ctx context.Context, app *amplifyx.App, cmd string, args amplifyx.Args) error {
	var err error

	if panicErr := oops.Recover(func() {
		switch cmd {
		case "deploy":
			err = app.Deploy(ctx, args.Deploy)
		case "version":
			err = app.Version(ctx, args.Version)
		default:
			err = oops.Errorf("unexpected command: %s", cmd)
		}
	}); panicErr != nil {
		err = panicErr
	}

	return err
}

func fail(err error) int {
	fmt.Fprintln(os.Stderr, err.Error())

	return 1
}
