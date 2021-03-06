package cmd

import (
	"context"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zostay/garotate/pkg/config"
	"go.uber.org/zap"
)

var (
	rootCmd *cobra.Command
	c       config.Config
	cfgFile string
	dryRun  bool
	verbose bool
	devMode bool
	logger  *zap.Logger
	ctx     context.Context
)

func init() {
	cobra.OnInitialize(initContext, initConfig)

	rootCmd = &cobra.Command{
		Use:   "garotate",
		Short: "tools for managing AWS secrets on github",
	}

	viper.SetDefault("RotateAfter", 168*time.Hour)
	viper.SetDefault("DisableAfter", 48*time.Hour)

	rootCmd.PersistentFlags().BoolVar(
		&devMode, "dev-mode", false,
		"turns on developer mode logging",
	)
	rootCmd.PersistentFlags().StringVar(
		&cfgFile, "config-file", "",
		"config file (default is /garotate.yaml)",
	)
	rootCmd.PersistentFlags().Duration(
		"rotate-after", 168*time.Hour,
		"keys older than rotate-after will be rotated",
	)
	rootCmd.PersistentFlags().Duration(
		"disable-after", 48*time.Hour,
		"keys older than rotate-after + disable-after will be disabled",
	)
	rootCmd.PersistentFlags().BoolVar(
		&dryRun, "dry-run", false,
		"a dry-run describes what would happen without doing it",
	)
	rootCmd.PersistentFlags().BoolVarP(
		&verbose, "verbose", "v", false,
		"more verbose logging",
	)

	err := viper.BindPFlag(
		"rotateAfter", rootCmd.PersistentFlags().Lookup("rotate-after"),
	)
	if err != nil {
		panic(err)
	}

	err = viper.BindPFlag(
		"disableAfter", rootCmd.PersistentFlags().Lookup("disable-after"),
	)
	if err != nil {
		panic(err)
	}

	viper.SetDefault("rotateAfter", 168*time.Hour)
	viper.SetDefault("disableAfter", 48*time.Hour)

	initRotateCmd()
	initDisableCmd()
}

func initContext() {
	if devMode {
		logger = config.DevelopmentLogger()
	} else {
		logger = config.ProductionLogger()
	}

	if verbose {
		logger = logger.WithOptions(
			zap.IncreaseLevel(zap.DebugLevel),
		)
	}

	ctx = config.WithLogger(context.Background(), logger)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("/")
		viper.SetConfigType("yaml")
		viper.SetConfigName("garotate")
	}

	viper.AutomaticEnv()

	slog := logger.Sugar()

	err := viper.ReadInConfig()
	if err != nil {
		slog.Fatalf("unable to read configuration: %v", err)
	}

	err = viper.Unmarshal(&c)
	if err != nil {
		slog.Fatalf("unable to unmarshal configuration: %v", err)
	}

	err = c.Prepare()
	if err != nil {
		slog.Fatalf("unable to finish processing configuration: %v", err)
	}
}

func Execute() {
	err := rootCmd.Execute()
	cobra.CheckErr(err)
}
