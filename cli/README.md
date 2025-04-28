# Sunbird CLI

A command-line interface for Sunbird data products, providing type-safe operations for running jobs, replaying data, and managing configurations.

## Features

- Type-safe Go implementation
- Run data processing jobs
- Replay jobs for specific date ranges
- Telemetry replay functionality
- Backup and restore capabilities

## Prerequisites

- Go 1.16 or higher
- Spark installation (set via SPARK_HOME environment variable)
- Project home directory (set via PROJECT_HOME environment variable)

## Installation

### From Source

```bash
git clone https://github.com/ChakshuGautam/sunbird-core-dataproducts.git
cd sunbird-core-dataproducts/cli
go build -o sunbird-cli
```

### Environment Variables

Set the following environment variables:

```bash
export SPARK_HOME=/path/to/spark
export PROJECT_HOME=/path/to/sunbird-core-dataproducts
export DP_LOGS=/path/to/logs  # Optional, defaults to $PROJECT_HOME/platform-scripts/shell/local/logs
```

## Usage

### Run a Job

```bash
./sunbird-cli run [model]
```

Example:
```bash
./sunbird-cli run monitor-job-summ
```

### Replay a Job

```bash
./sunbird-cli replay [model] [start-date] [end-date]
```

Example:
```bash
./sunbird-cli replay wfs 2023-01-01 2023-01-31
```

### Replay an Updater Job

```bash
./sunbird-cli replay-updater [model] [start-date] [end-date]
```

Example:
```bash
./sunbird-cli replay-updater wfu 2023-01-01 2023-01-31
```

### Replay Telemetry

```bash
./sunbird-cli telemetry-replay [end-date]
```

Example:
```bash
./sunbird-cli telemetry-replay 2023-01-31
```

## Available Models

- `monitor-job-summ`: Monitor job summary
- `job-manager`: Job manager
- `wfs`: Workflow summary
- `wfus`: Workflow usage summary
- `wfu`: Workflow usage updater
- `telemetry-replay`: Telemetry replay
- `assessment-dashboard-metrics`: Assessment dashboard metrics

## License

This project is licensed under the same license as the Sunbird Core Data Products.
