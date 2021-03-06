package cmd

import (
	"github.com/spf13/cobra"

	"github.com/zostay/garotate/pkg/config"
	"github.com/zostay/garotate/pkg/plugin"
	"github.com/zostay/garotate/pkg/rotate"
)

var (
	rotateCmd   *cobra.Command
	alsoDisable bool
)

// initRotateCmd configures the command.
func initRotateCmd() {
	rotateCmd = &cobra.Command{
		Use:   "rotate",
		Short: "rotate an AWS IAM key/secret pair and update a github action secret",
		Run:   RunRotation,
	}

	rotateCmd.Flags().BoolVar(&alsoDisable, "also-disable", false, "after rotating keys, check for any old access keys that should be disabled")

	rootCmd.AddCommand(rotateCmd)
}

// RunRotation performs all the rotations configured according to the command
// line given. It may also trigger disablement, if the --also-disable option is
// set.
func RunRotation(cmd *cobra.Command, args []string) {
	buildMgr := plugin.NewManager(c.Plugins)
	for _, r := range c.Rotations {
		RunRotations(buildMgr, &r)
	}
	if alsoDisable {
		for _, r := range c.Disablements {
			RunDisablement(buildMgr, &r)
		}
	}
}

// RunRotations performs a single rotation from the configuration.
func RunRotations(
	buildMgr *plugin.Manager,
	r *config.Rotation,
) {
	slog := logger.Sugar()

	rc, err := buildMgr.Instance(ctx, r.RotateClient)
	rotCli, ok := rc.(rotate.Client)
	if !ok {
		slog.Errorw(
			"failed to load rotation client",
			"client_name", r.RotateClient,
			"error", err,
		)
		return
	}

	secretSet, err := findSecretSet(r.SecretSet)
	if err != nil {
		slog.Errorw(
			"failed to locate the secret set to work with ",
			"client_name", r.RotateClient,
			"client_desc", rotCli.Name(),
			"error", err,
		)
		return
	}

	m := rotate.New(
		rotCli,
		r.RotateAfter,
		dryRun,
		buildMgr,
		secretSet.Secrets,
	)

	err = m.RotateSecrets(ctx)
	if err != nil {
		slog.Errorw(
			"failed to complete secret rotation",
			"client_name", r.RotateClient,
			"client_desc", rotCli.Name(),
			"error", err,
		)
	}
}
