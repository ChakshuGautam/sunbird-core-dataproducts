# Sunbird Core Data Products

Repository for generic data products to derive reports and insights from telemetry. The repository consists of base interfaces and implementations for the entire data product life cycle.

## Components

### Batch Models

Batch processing models for data analysis and transformation.

### Job Manager

Manages the execution of data processing jobs.

### Video Streaming

Components for video streaming analytics.

### CLI (Go)

A type-safe command-line interface for Sunbird data products, providing operations for running jobs, replaying data, and managing configurations.

#### Features

- Type-safe Go implementation
- Run data processing jobs
- Replay jobs for specific date ranges
- Telemetry replay functionality
- Backup and restore capabilities

For more information, see the [CLI README](cli/README.md).

## Getting Started

### Prerequisites

- Java 8
- Scala 2.11
- Apache Spark 2.0.1
- Go 1.16+ (for CLI)

### Building

```bash
mvn clean install
```

### Building the CLI

```bash
cd cli
go build -o sunbird-cli
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.
