package utils

import (
	"fmt"
	"os/exec"
	"time"
)

// BackupManager handles backup and restore operations for data
type BackupManager struct {
	AwsCmd string
}

// NewBackupManager creates a new BackupManager
func NewBackupManager(awsCmd string) *BackupManager {
	if awsCmd == "" {
		awsCmd = "aws"
	}
	return &BackupManager{
		AwsCmd: awsCmd,
	}
}

// Backup creates a backup of data for a specific date range
func (b *BackupManager) Backup(startDate, endDate, bucketName, prefix, backupDir string) error {
	startTime, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return fmt.Errorf("invalid start date format: %w", err)
	}

	endTime, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return fmt.Errorf("invalid end date format: %w", err)
	}

	src := fmt.Sprintf("s3://%s/%s/", bucketName, prefix)
	dst := fmt.Sprintf("s3://%s/%s/", bucketName, backupDir)

	fmt.Printf("Backing up the files from %s to %s for the date range - (%s, %s)\n", src, dst, startDate, endDate)

	// Loop through each day in the date range
	for currentTime := startTime; !currentTime.After(endTime); currentTime = currentTime.AddDate(0, 0, 1) {
		date := currentTime.Format("2006-01-02")
		
		// AWS S3 move command
		cmd := exec.Command(
			b.AwsCmd, "s3", "mv", src, dst,
			"--recursive",
			"--exclude", "*",
			"--include", fmt.Sprintf("%s-*", date),
		)

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to backup files for date %s: %w", date, err)
		}
	}

	return nil
}

// Rollback restores data from a backup
func (b *BackupManager) Rollback(bucketName, prefix, backupDir string) error {
	src := fmt.Sprintf("s3://%s/%s/", bucketName, backupDir)
	dst := fmt.Sprintf("s3://%s/%s/", bucketName, prefix)

	fmt.Printf("Copy back the %s files to source directory %s from backup directory %s\n", prefix, dst, src)

	// AWS S3 copy command
	cmd := exec.Command(
		b.AwsCmd, "s3", "cp", src, dst,
		"--recursive",
		"--include", "*",
	)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to rollback files: %w", err)
	}

	return nil
}

// Delete removes a backup directory
func (b *BackupManager) Delete(bucketName, backupDir string) error {
	path := fmt.Sprintf("s3://%s/%s/", bucketName, backupDir)
	fmt.Printf("Deleting the back-up files from %s\n", path)

	// AWS S3 remove command
	cmd := exec.Command(
		b.AwsCmd, "s3", "rm", path,
		"--recursive",
	)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to delete backup files: %w", err)
	}

	return nil
}
