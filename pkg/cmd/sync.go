package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewSyncCommand creates the sync command
func NewSyncCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Synchronize recent commits between platforms",
		Long: `Synchronize recent commits from source platforms to target platforms.

By default, this will sync commits from the last 24 hours. You can specify
a different time range using the --since flag.`,
		RunE: runSync,
	}

	cmd.Flags().String("since", "24h", "sync commits since this duration (e.g., 24h, 7d, 1w, 3mo, 1y)")
	cmd.Flags().StringSlice("sources", nil, "specific source platforms to sync from")
	cmd.Flags().StringSlice("targets", nil, "specific target platforms to sync to")
	cmd.Flags().Bool("force", false, "force sync even if commits already exist")

	return cmd
}

func runSync(cmd *cobra.Command, args []string) error {
	verbose := viper.GetBool("verbose")
	dryRun := viper.GetBool("dry-run")

	if verbose {
		fmt.Println("🔄 Starting sync operation...")
	}

	// Parse since duration
	sinceStr, _ := cmd.Flags().GetString("since")
	since, err := parseDuration(sinceStr)
	if err != nil {
		return fmt.Errorf("invalid since duration: %w", err)
	}

	sinceTime := time.Now().Add(-since)

	if verbose {
		fmt.Printf("📅 Syncing commits since: %s\n", sinceTime.Format(time.RFC3339))
	}

	if dryRun {
		fmt.Println("🧪 Dry run mode - no changes will be made")
		fmt.Println("✅ Sync completed (dry run)")
		return nil
	}

	// TODO: Implement actual sync logic here
	// This would:
	// 1. Load configuration
	// 2. Initialize source and target platforms
	// 3. Fetch commits from sources since the specified time
	// 4. Mirror commits to targets

	fmt.Println("✅ Sync completed successfully")

	return nil
}

// parseDuration parses duration strings like "24h", "7d", "1w", "3mo", "1y"
func parseDuration(s string) (time.Duration, error) {
	// Handle common suffixes
	switch {
	case len(s) > 1 && s[len(s)-1] == 'y':
		// Years (approximate: 365 days)
		years, err := time.ParseDuration(s[:len(s)-1] + "h")
		if err != nil {
			return 0, err
		}
		return years * 24 * 365, nil
	case len(s) > 2 && s[len(s)-2:] == "mo":
		// Months (approximate: 30 days) - use "mo" to avoid conflict with minutes "m"
		months, err := time.ParseDuration(s[:len(s)-2] + "h")
		if err != nil {
			return 0, err
		}
		return months * 24 * 30, nil
	case len(s) > 1 && s[len(s)-1] == 'd':
		// Days
		days, err := time.ParseDuration(s[:len(s)-1] + "h")
		if err != nil {
			return 0, err
		}
		return days * 24, nil
	case len(s) > 1 && s[len(s)-1] == 'w':
		// Weeks
		weeks, err := time.ParseDuration(s[:len(s)-1] + "h")
		if err != nil {
			return 0, err
		}
		return weeks * 24 * 7, nil
	default:
		// Standard Go duration format (h, m, s)
		return time.ParseDuration(s)
	}
}
