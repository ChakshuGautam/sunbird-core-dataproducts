package cmd

import (
	"fmt"
	"os"

	"github.com/ChakshuGautam/sunbird-core-dataproducts/cli/internal/config"
	"github.com/ChakshuGautam/sunbird-core-dataproducts/cli/internal/job"
	"github.com/spf13/cobra"
)

var telemetryReplayCmd = &cobra.Command{
	Use:   "telemetry-replay [end-date]",
	Short: "Replay telemetry events",
	Long: `Replay telemetry events for a specific date.
Example: sunbird-cli telemetry-replay 2023-01-31`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		endDate := args[0]

		// Get environment variables
		sparkHome := os.Getenv("SPARK_HOME")
		if sparkHome == "" {
			er("SPARK_HOME environment variable is not set")
		}

		projectHome := os.Getenv("PROJECT_HOME")
		if projectHome == "" {
			er("PROJECT_HOME environment variable is not set")
		}

		logDir := os.Getenv("DP_LOGS")
		if logDir == "" {
			logDir = fmt.Sprintf("%s/platform-scripts/shell/local/logs", projectHome)
		}

		// Get model config
		cfg := config.NewConfig()
		configStr, err := cfg.GetModelConfig("telemetry-replay", endDate)
		if err != nil {
			er(fmt.Errorf("failed to get model config: %w", err))
		}

		// Create job executor
		executor := job.NewJobExecutor(sparkHome, projectHome, logDir)

		// Run the job
		if err := executor.RunJob("telemetry-replay", configStr); err != nil {
			er(fmt.Errorf("telemetry replay failed: %w", err))
		}

		fmt.Println("Telemetry replay executed successfully")
	},
}

func init() {
	rootCmd.AddCommand(telemetryReplayCmd)
}
