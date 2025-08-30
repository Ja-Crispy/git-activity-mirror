package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewRootCommand creates the root cobra command
func NewRootCommand(version string) *cobra.Command {
	var cfgFile string

	rootCmd := &cobra.Command{
		Use:   "git-activity-mirror",
		Short: "Mirror git commit activity between platforms",
		Long: `git-activity-mirror is a cross-platform tool that mirrors git commit activity 
between any git hosting platforms (GitHub, GitLab, Bitbucket, etc.) while maintaining privacy.

Your actual code is never exposed - only commit timestamps and generic messages are mirrored,
allowing you to maintain an accurate contribution graph across all platforms.`,
		Version: version,
	}

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.git-activity-mirror/config.yaml)")
	rootCmd.PersistentFlags().Bool("verbose", false, "verbose output")
	rootCmd.PersistentFlags().Bool("dry-run", false, "show what would be done without making changes")

	// Bind flags to viper
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("dry-run", rootCmd.PersistentFlags().Lookup("dry-run"))

	// Initialize config
	cobra.OnInitialize(func() { initConfig(cfgFile) })

	// Add subcommands
	rootCmd.AddCommand(NewInitCommand())
	rootCmd.AddCommand(NewSyncCommand())
	rootCmd.AddCommand(NewImportCommand())
	rootCmd.AddCommand(NewStatusCommand())
	rootCmd.AddCommand(NewConfigCommand())

	return rootCmd
}

// initConfig reads in config file and ENV variables
func initConfig(cfgFile string) {
	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".git-activity-mirror"
		configDir := filepath.Join(home, ".git-activity-mirror")
		viper.AddConfigPath(configDir)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in
	if err := viper.ReadInConfig(); err == nil {
		if viper.GetBool("verbose") {
			fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		}
	}
}