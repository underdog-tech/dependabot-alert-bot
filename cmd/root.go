package cmd

import (
	"path/filepath"

	"github.com/rs/zerolog"
	"github.com/underdog-tech/vulnbot/internal"
	"github.com/underdog-tech/vulnbot/logger"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "vulnbot",
	Short: "Vulnbot: Your all-in-one security alert manager.",
	Long: `Vulnbot is a comprehensive security alert manager designed to keep your code safe from vulnerabilities.

It is a versatile bot that can seamlessly integrate with multiple data sources, such as GitHub, and soon Phylum,
Vulnbot empowers developers and security teams to efficiently manage and respond to security threats.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		quiet, _ := cmd.Flags().GetBool("quiet")
		verbosity, _ := cmd.Flags().GetCount("verbose")

		if quiet == true {
			logger.SetLogLevel(zerolog.Disabled)
		} else if verbosity > 0 {
			if verbosity > 3 {
				verbosity = 3
			}
			logLevel := logger.DEFAULT_LOG_LEVEL - zerolog.Level(verbosity)
			logger.SetLogLevel(logLevel)
		} else {
			logger.SetLogLevel(logger.DEFAULT_LOG_LEVEL)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	log := logger.Get()
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to execute command.")
	}
}

func init() {
	persistent := rootCmd.PersistentFlags()
	persistent.BoolP("disable-slack", "d", false, "Disable Slack alerts.")

	projectRootDir := internal.GetProjectRootDir()
	persistent.StringP("config", "c", filepath.Join(projectRootDir, "config.toml"), "Config file path.")

	persistent.BoolP("quiet", "q", false, "Suppress all console output. (Mutually exclusive with 'verbose'.)")
	persistent.CountP("verbose", "v", "More verbose output. Specifying multiple times increases verbosity. (Mutually exclusive with 'quiet'.)")
	rootCmd.MarkFlagsMutuallyExclusive("verbose", "quiet")
}
