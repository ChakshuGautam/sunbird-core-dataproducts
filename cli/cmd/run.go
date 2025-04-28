package cmd

import (
	"fmt"
	"os"

	"github.com/ChakshuGautam/sunbird-core-dataproducts/cli/internal/config"
	"github.com/ChakshuGautam/sunbird-core-dataproducts/cli/internal/job"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [model]",
	Short: "Run a data processing job",
	Long: `Run a data processing job with the specified model.
Example: sunbird-cli run monitor-job-summ`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		model := args[0]
		configStr, err := cmd.Flags().GetString("config")
		if err != nil {
			er(fmt.Errorf("failed to get config flag: %w", err))
		}

		// If config is not provided, get it from the model config
		if configStr == "" {
			cfg := config.NewConfig()
			endDate, _ := cmd.Flags().GetString("end-date")
			configStr, err = cfg.GetModelConfig(model, endDate)
			if err != nil {
				er(fmt.Errorf("failed to get model config: %w", err))
			}
		}

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

		// Create job executor
		executor := job.NewJobExecutor(sparkHome, projectHome, logDir)

		// Run the job
		if err := executor.RunJob(model, configStr); err != nil {
			er(fmt.Errorf("job execution failed: %w", err))
		}

		fmt.Printf("Job %s executed successfully\n", model)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Add flags
	runCmd.Flags().StringP("config", "c", "", "Custom configuration for the job (JSON format)")
	runCmd.Flags().StringP("end-date", "e", "", "End date for the job (format: YYYY-MM-DD)")
}
