package config

import (
	"fmt"
	"time"
)

// ModelConfig represents the configuration for a data model
type ModelConfig struct {
	Search       Search       `json:"search"`
	Model        string       `json:"model"`
	ModelParams  ModelParams  `json:"modelParams,omitempty"`
	Output       []Output     `json:"output"`
	Parallelization int       `json:"parallelization,omitempty"`
	AppName      string       `json:"appName"`
	DeviceMapping bool        `json:"deviceMapping"`
}

// Search represents the search configuration
type Search struct {
	Type    string   `json:"type"`
	Queries []Query  `json:"queries,omitempty"`
}

// Query represents a search query
type Query struct {
	Bucket   string `json:"bucket,omitempty"`
	Prefix   string `json:"prefix,omitempty"`
	EndDate  string `json:"endDate,omitempty"`
	Delta    int    `json:"delta,omitempty"`
	File     string `json:"file,omitempty"`
}

// ModelParams represents parameters for the model
type ModelParams struct {
	ApiVersion     string   `json:"apiVersion,omitempty"`
	PushMetrics    bool     `json:"pushMetrics,omitempty"`
	BrokerList     string   `json:"brokerList,omitempty"`
	Topic          string   `json:"topic,omitempty"`
	Model          []ModelInfo `json:"model,omitempty"`
	FromDate       string   `json:"fromDate,omitempty"`
	ToDate         string   `json:"toDate,omitempty"`
}

// ModelInfo represents information about a model
type ModelInfo struct {
	Model           string `json:"model"`
	Category        string `json:"category"`
	InputDependency string `json:"input_dependency"`
}

// Output represents an output configuration
type Output struct {
	To      string       `json:"to"`
	Params  OutputParams `json:"params"`
}

// OutputParams represents parameters for output
type OutputParams struct {
	PrintEvent bool   `json:"printEvent,omitempty"`
	BrokerList string `json:"brokerList,omitempty"`
	Topic      string `json:"topic,omitempty"`
}

// JobManagerConfig represents the configuration for the job manager
type JobManagerConfig struct {
	JobsCount        int    `json:"jobsCount"`
	Topic            string `json:"topic"`
	BootStrapServer  string `json:"bootStrapServer"`
	ZookeeperConnect string `json:"zookeeperConnect"`
	ConsumerGroup    string `json:"consumerGroup"`
	SlackChannel     string `json:"slackChannel"`
	SlackUserName    string `json:"slackUserName"`
	TempBucket       string `json:"tempBucket"`
	TempFolder       string `json:"tempFolder"`
}

// Config holds the environment configuration
type Config struct {
	DatasetRawBucket            string
	DataExhaustBucket           string
	DataExhaustPrefix           string
	ConsumptionRawPrefix        string
	JobTopic                    string
	Topic                       string
	AnalyticsMetricsTopic       string
	LearningTopic               string
	CurrentDate                 string
	AnalyticsHome               string
	TempFolder                  string
	Bucket                      string
	Zookeeper                   string
	BrokerList                  string
	InputBucket                 string
	SinkTopic                   string
}

// NewConfig creates a new configuration with default values
func NewConfig() *Config {
	currentDate := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	return &Config{
		DatasetRawBucket:      "ekstep-data-sets-dev",
		DataExhaustBucket:     "ekstep-public-dev",
		DataExhaustPrefix:     "data-exhaust/",
		ConsumptionRawPrefix:  "datasets/D001/4208ab995984d222b59299e5103d350a842d8d41/",
		JobTopic:              "analytics.job_queue",
		Topic:                 "job-manager-test",
		AnalyticsMetricsTopic: "local.telemetry.derived",
		LearningTopic:         "local.telemetry.derived",
		CurrentDate:           currentDate,
		AnalyticsHome:         "/home",
		TempFolder:            "transient-data",
		Bucket:                "dev-data-store",
		Zookeeper:             "localhost:2181",
		BrokerList:            "localhost:9092",
		InputBucket:           "extractor-failed/",
		SinkTopic:             "sunbirddev.telemetry.ingest.replay",
	}
}

// GetModelConfig returns the configuration for a specific model
func (c *Config) GetModelConfig(modelCode string, endDate string) (string, error) {
	if endDate == "" {
		endDate = time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	}

	switch modelCode {
	case "monitor-job-summ":
		return c.getMonitorJobSummConfig(), nil
	case "job-manager":
		return c.getJobManagerConfig(), nil
	case "wfs":
		return c.getWorkflowSummaryConfig(), nil
	case "wfus":
		return c.getWorkflowUsageSummaryConfig(endDate), nil
	case "wfu":
		return c.getWorkflowUsageUpdaterConfig(endDate), nil
	case "telemetry-replay":
		return c.getTelemetryReplayConfig(endDate), nil
	case "assessment-dashboard-metrics":
		return c.getAssessmentDashboardMetricsConfig(), nil
	default:
		return "", fmt.Errorf("unknown model code: %s", modelCode)
	}
}

// Implementation of specific model configurations
func (c *Config) getMonitorJobSummConfig() string {
	// This would return the JSON configuration for monitor-job-summ
	// For brevity, I'm returning a placeholder
	return fmt.Sprintf(`{
		"search": {
			"type": "local",
			"queries": [{
				"file": "%s/scripts/logs/joblog.log"
			}]
		},
		"model": "org.ekstep.analytics.model.MonitorSummaryModel",
		"modelParams": {
			"pushMetrics": true,
			"brokerList": "%s",
			"topic": "%s",
			"model": []
		},
		"output": [{
			"to": "console",
			"params": {
				"printEvent": false
			}
		}, {
			"to": "kafka",
			"params": {
				"brokerList": "%s",
				"topic": "%s"
			}
		}],
		"appName": "TestMonitorSummarizer",
		"deviceMapping": true
	}`, c.AnalyticsHome, c.BrokerList, c.AnalyticsMetricsTopic, c.BrokerList, c.Topic)
}

func (c *Config) getJobManagerConfig() string {
	return fmt.Sprintf(`{
		"jobsCount": 1,
		"topic": "%s",
		"bootStrapServer": "%s",
		"zookeeperConnect": "%s",
		"consumerGroup": "jobmanager",
		"slackChannel": "#testing",
		"slackUserName": "JobManager",
		"tempBucket": "%s",
		"tempFolder": "%s"
	}`, c.JobTopic, c.BrokerList, c.Zookeeper, c.Bucket, c.TempFolder)
}

func (c *Config) getWorkflowSummaryConfig() string {
	return fmt.Sprintf(`{
		"search": {
			"type": "azure",
			"queries": [{
				"bucket": "%s",
				"prefix": "unique/",
				"endDate": "2019-09-23",
				"delta": 0
			}]
		},
		"model": "org.ekstep.analytics.model.WorkflowSummary",
		"modelParams": {
			"apiVersion": "v2"
		},
		"output": [{
			"to": "console",
			"params": {
				"printEvent": false
			}
		}],
		"parallelization": 8,
		"appName": "Workflow Summarizer",
		"deviceMapping": true
	}`, c.Bucket)
}

func (c *Config) getWorkflowUsageSummaryConfig(endDate string) string {
	return fmt.Sprintf(`{
		"search": {
			"type": "s3",
			"queries": [{
				"bucket": "%s",
				"prefix": "wfs/",
				"endDate": "%s",
				"delta": 0
			}]
		},
		"model": "org.ekstep.analytics.model.WorkflowUsageSummary",
		"modelParams": {
			"apiVersion": "v2"
		},
		"output": [{
			"to": "console",
			"params": {
				"printEvent": false
			}
		}, {
			"to": "kafka",
			"params": {
				"brokerList": "%s",
				"topic": "%s"
			}
		}],
		"parallelization": 8,
		"appName": "Workflow Usage Summarizer",
		"deviceMapping": false
	}`, c.Bucket, endDate, c.BrokerList, c.Topic)
}

func (c *Config) getWorkflowUsageUpdaterConfig(endDate string) string {
	return fmt.Sprintf(`{
		"search": {
			"type": "s3",
			"queries": [{
				"bucket": "%s",
				"prefix": "wfus/",
				"endDate": "%s",
				"delta": 0
			}]
		},
		"model": "org.ekstep.analytics.updater.UpdateWorkFlowUsageDB",
		"output": [{
			"to": "console",
			"params": {
				"printEvent": false
			}
		}],
		"parallelization": 10,
		"appName": "Workflow Usage Updater",
		"deviceMapping": false
	}`, c.Bucket, endDate)
}

func (c *Config) getTelemetryReplayConfig(endDate string) string {
	return fmt.Sprintf(`{
		"search": {
			"type": "azure",
			"queries": [{
				"bucket": "%s",
				"prefix": "%s",
				"endDate": "%s",
				"delta": 0
			}]
		},
		"model": "org.ekstep.analytics.job.EventsReplayJob",
		"modelParams": {},
		"output": [{
			"to": "console",
			"params": {
				"printEvent": false
			}
		}, {
			"to": "kafka",
			"params": {
				"brokerList": "%s",
				"topic": "%s"
			}
		}],
		"parallelization": 8,
		"appName": "TelemetryReplayJob",
		"deviceMapping": false
	}`, c.Bucket, c.InputBucket, endDate, c.BrokerList, c.SinkTopic)
}

func (c *Config) getAssessmentDashboardMetricsConfig() string {
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	return fmt.Sprintf(`{
		"search": {
			"type": "none"
		},
		"model": "org.ekstep.analytics.job.AssessmentMetricsJob",
		"modelParams": {
			"fromDate": "%s",
			"toDate": "%s",
			"sparkCassandraConnectionHost": "11.2.3.63",
			"sparkElasticsearchConnectionHost": "11.2.3.58"
		},
		"output": [{
			"to": "console",
			"params": {
				"printEvent": false
			}
		}],
		"parallelization": 8,
		"appName": "Assessment Dashboard Metrics",
		"deviceMapping": false
	}`, yesterday, yesterday)
}
