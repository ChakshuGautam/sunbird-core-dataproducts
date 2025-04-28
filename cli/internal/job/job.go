package job

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// JobExecutor handles the execution of Spark jobs
type JobExecutor struct {
	SparkHome    string
	ProjectHome  string
	LogDirectory string
}

// NewJobExecutor creates a new JobExecutor with the given configuration
func NewJobExecutor(sparkHome, projectHome, logDirectory string) *JobExecutor {
	return &JobExecutor{
		SparkHome:    sparkHome,
		ProjectHome:  projectHome,
		LogDirectory: logDirectory,
	}
}

// RunJob executes a job with the given model and configuration
func (j *JobExecutor) RunJob(model, config string) error {
	today := time.Now().Format("2006-01-02")
	logFile := filepath.Join(j.LogDirectory, fmt.Sprintf("%s-job-execution.log", today))

	// Ensure log directory exists
	if err := os.MkdirAll(j.LogDirectory, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	// Log job start
	logEntry := fmt.Sprintf("Starting the job - %s\n", model)
	if err := appendToFile(logFile, logEntry); err != nil {
		return fmt.Errorf("failed to write to log file: %w", err)
	}

	// Prepare the spark-submit command
	cmd := exec.Command(
		filepath.Join(j.SparkHome, "bin", "spark-submit"),
		"--master", "local[*]",
		"--jars", filepath.Join(j.ProjectHome, "platform-framework/analytics-job-driver/target/analytics-framework-1.0.jar"),
		"--class", "org.ekstep.analytics.job.JobExecutor",
		filepath.Join(j.ProjectHome, "platform-modules/batch-models/target/batch-models-1.0.jar"),
		"--model", model,
		"--config", config,
	)

	// Redirect output to log file
	logFileHandle, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer logFileHandle.Close()

	cmd.Stdout = logFileHandle
	cmd.Stderr = logFileHandle

	// Execute the command
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("job execution failed: %w", err)
	}

	// Log job completion
	logEntry = fmt.Sprintf("Job execution completed - %s\n", model)
	if err := appendToFile(logFile, logEntry); err != nil {
		return fmt.Errorf("failed to write to log file: %w", err)
	}

	return nil
}

// ReplayJob replays a job for a specific date range
func (j *JobExecutor) ReplayJob(model, config, startDate, endDate string) error {
	// Ensure log directory exists
	if err := os.MkdirAll(j.LogDirectory, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	logFile := filepath.Join(j.LogDirectory, fmt.Sprintf("%s-%s-replay.log", endDate, model))

	// Prepare the spark-submit command for replay
	cmd := exec.Command(
		filepath.Join(j.SparkHome, "bin", "spark-submit"),
		"--master", "local[*]",
		"--jars", filepath.Join(j.ProjectHome, "platform-framework/analytics-job-driver/target/analytics-framework-1.0.jar"),
		"--class", "org.ekstep.analytics.job.ReplaySupervisor",
		filepath.Join(j.ProjectHome, "platform-modules/batch-models/target/batch-models-1.0.jar"),
		"--model", model,
		"--fromDate", startDate,
		"--toDate", endDate,
		"--config", config,
	)

	// Redirect output to log file
	logFileHandle, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer logFileHandle.Close()

	cmd.Stdout = logFileHandle
	cmd.Stderr = logFileHandle

	// Execute the command
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("job replay failed: %w", err)
	}

	return nil
}

// appendToFile appends text to a file
func appendToFile(filename, text string) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.WriteString(text); err != nil {
		return err
	}
	return nil
}
