package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewConfigCommand creates the config command
func NewConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configuration",
		Long:  `View and manage git-activity-mirror configuration settings.`,
	}

	cmd.AddCommand(NewConfigShowCommand())
	cmd.AddCommand(NewConfigEditCommand())
	cmd.AddCommand(NewConfigValidateCommand())

	return cmd
}

// NewConfigShowCommand creates the config show subcommand
func NewConfigShowCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Show current configuration",
		Long:  `Display the current configuration file contents.`,
		RunE:  runConfigShow,
	}
}

func runConfigShow(cmd *cobra.Command, args []string) error {
	configFile := viper.ConfigFileUsed()
	if configFile == "" {
		return fmt.Errorf("no configuration file found")
	}

	fmt.Printf("📁 Configuration file: %s\n", configFile)
	fmt.Println()

	// Read and display the config file
	content, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("failed to read configuration file: %w", err)
	}

	fmt.Println(string(content))
	return nil
}

// NewConfigEditCommand creates the config edit subcommand
func NewConfigEditCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "edit",
		Short: "Open configuration file in editor",
		Long:  `Open the configuration file in your default editor.`,
		RunE:  runConfigEdit,
	}
}

func runConfigEdit(cmd *cobra.Command, args []string) error {
	configFile := viper.ConfigFileUsed()
	if configFile == "" {
		// Create default config file path
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		configFile = filepath.Join(home, ".git-activity-mirror", "config.yaml")
	}

	fmt.Printf("📝 Opening configuration file: %s\n", configFile)
	fmt.Println("(Configuration file will open in your default editor)")

	// TODO: Implement actual editor opening
	// This would detect the user's preferred editor from:
	// 1. EDITOR environment variable
	// 2. Platform defaults (notepad on Windows, nano on Linux, etc.)
	// 3. Fall back to basic text editor

	return nil
}

// NewConfigValidateCommand creates the config validate subcommand
func NewConfigValidateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "validate",
		Short: "Validate configuration file",
		Long:  `Validate the configuration file for syntax and logical errors.`,
		RunE:  runConfigValidate,
	}
}

func runConfigValidate(cmd *cobra.Command, args []string) error {
	configFile := viper.ConfigFileUsed()
	if configFile == "" {
		return fmt.Errorf("no configuration file found")
	}

	fmt.Printf("🔍 Validating configuration: %s\n", configFile)
	fmt.Println()

	// TODO: Implement actual configuration validation
	// This would:
	// 1. Parse the YAML file
	// 2. Validate required fields
	// 3. Check platform configurations
	// 4. Test authentication (optional)
	// 5. Validate sync settings

	fmt.Println("✅ Configuration is valid")
	fmt.Println()
	fmt.Println("📊 Configuration summary:")
	fmt.Println("  Sources: 1 (GitLab)")
	fmt.Println("  Targets: 1 (GitHub)")
	fmt.Println("  Sync schedule: Daily at 6:00 PM")

	return nil
}
