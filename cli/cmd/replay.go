package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/ChakshuGautam/sunbird-core-dataproducts/cli/internal/config"
	"github.com/ChakshuGautam/sunbird-core-dataproducts/cli/internal/job"
	"github.com/ChakshuGautam/sunbird-core-dataproducts/cli/internal/utils"
	"github.com/spf13/cobra"
)

var replayCmd = &cobra.Command{
	Use:   "replay [model] [start-date] [end-date]",
	Short: "Replay a data processing job for a specific date range",
	Long: `Replay a data processing job with the specified model for a date range.
Example: sunbird-cli replay wfs 2023-01-01 2023-01-31`,
	Args: cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		model := args[0]
		startDate := args[1]
		endDate := args[2]

		// Validate dates
		_, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			er(fmt.Errorf("invalid start date format (should be YYYY-MM-DD): %w", err))
		}

		_, err = time.Parse("2006-01-02", endDate)
		if err != nil {
			er(fmt.Errorf("invalid end date format (should be YYYY-MM-DD): %w", err))
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

		// Get model config
		cfg := config.NewConfig()
		configStr, err := cfg.GetModelConfig(model, endDate)
		if err != nil {
			er(fmt.Errorf("failed to get model config: %w", err))
		}

		// Create backup manager
		backupManager := utils.NewBackupManager("")

		// Create backup
		bucketName := "sandbox-data-store"
		backupDir := fmt.Sprintf("backup-%s", model)

		fmt.Printf("Creating backup for %s from %s to %s...\n", model, startDate, endDate)
		err = backupManager.Backup(startDate, endDate, bucketName, model, backupDir)
		if err != nil {
			er(fmt.Errorf("backup failed: %w", err))
		}
		fmt.Println("Backup completed successfully")

		// Create job executor
		executor := job.NewJobExecutor(sparkHome, projectHome, logDir)

		// Run the replay job
		fmt.Printf("Running the %s job replay...\n", model)
		err = executor.ReplayJob(model, configStr, startDate, endDate)
		if err != nil {
			// If replay fails, rollback from backup
			fmt.Printf("Replay failed, rolling back from backup...\n")
			rollbackErr := backupManager.Rollback(bucketName, model, backupDir)
			if rollbackErr != nil {
				er(fmt.Errorf("rollback failed: %w", rollbackErr))
			}
			
			// Delete backup
			deleteErr := backupManager.Delete(bucketName, backupDir)
			if deleteErr != nil {
				fmt.Printf("Warning: Failed to delete backup: %v\n", deleteErr)
			}
			
			er(fmt.Errorf("job replay failed: %w", err))
		}

		// Delete backup after successful replay
		fmt.Println("Replay completed successfully, cleaning up backup...")
		err = backupManager.Delete(bucketName, backupDir)
		if err != nil {
			fmt.Printf("Warning: Failed to delete backup: %v\n", err)
		}

		fmt.Printf("Job %s replay executed successfully\n", model)
	},
}

var replayUpdaterCmd = &cobra.Command{
	Use:   "replay-updater [model] [start-date] [end-date]",
	Short: "Replay an updater job for a specific date range",
	Long: `Replay an updater job with the specified model for a date range.
Example: sunbird-cli replay-updater wfu 2023-01-01 2023-01-31`,
	Args: cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		model := args[0]
		startDate := args[1]
		endDate := args[2]

		// Validate dates
		_, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			er(fmt.Errorf("invalid start date format (should be YYYY-MM-DD): %w", err))
		}

		_, err = time.Parse("2006-01-02", endDate)
		if err != nil {
			er(fmt.Errorf("invalid end date format (should be YYYY-MM-DD): %w", err))
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

		// Get model config
		cfg := config.NewConfig()
		configStr, err := cfg.GetModelConfig(model, endDate)
		if err != nil {
			er(fmt.Errorf("failed to get model config: %w", err))
		}

		// Create job executor
		executor := job.NewJobExecutor(sparkHome, projectHome, logDir)

		// Run the replay updater job
		fmt.Printf("Running the %s updater replay...\n", model)
		err = executor.ReplayJob(model, configStr, startDate, endDate)
		if err != nil {
			er(fmt.Errorf("updater replay failed: %w", err))
		}

		fmt.Printf("Job %s updater replay executed successfully\n", model)
	},
}

func init() {
	rootCmd.AddCommand(replayCmd)
	rootCmd.AddCommand(replayUpdaterCmd)
}
