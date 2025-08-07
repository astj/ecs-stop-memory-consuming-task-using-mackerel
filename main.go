package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mackerelio/mackerel-client-go"
)

type Config struct {
	DryRun          bool
	Verbose         bool
	MackerelService string
	MackerelRole    string
	MackerelMetric  string
	AWSProfile      string
	AWSRegion       string
	MackerelAPIKey  string
}

func main() {
	config := parseFlags()

	client := mackerel.NewClient(config.MackerelAPIKey)
	client.Verbose = config.Verbose

	arn, err := FindMostMemoryConsumingTaskArn(client, config.MackerelService, config.MackerelRole, config.MackerelMetric)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding most memory consuming task: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Most memory consuming task ARN:", arn)
}

func parseFlags() *Config {
	config := &Config{}

	// コマンドラインフラグの定義
	flag.BoolVar(&config.DryRun, "dry-run", false, "Dry run mode (no actual task termination)")
	flag.StringVar(&config.MackerelService, "mackerel-service", "", "Mackerel service name")
	flag.StringVar(&config.MackerelRole, "mackerel-role", "", "Mackerel role name")
	flag.StringVar(&config.MackerelMetric, "mackerel-metric", "", "Mackerel metric name for memory consumption")
	flag.StringVar(&config.AWSProfile, "aws-profile", "", "AWS profile name")
	flag.StringVar(&config.AWSRegion, "aws-region", "", "AWS region")
	flag.StringVar(&config.MackerelAPIKey, "mackerel-api-key", "", "Mackerel API key")
	flag.BoolVar(&config.Verbose, "verbose", false, "Enable verbose output")

	flag.Parse()

	// 環境変数からの設定（コマンドラインオプションが優先）
	if config.DryRun == false {
		if os.Getenv("DRY_RUN") == "true" {
			config.DryRun = true
		}
	}

	if config.MackerelService == "" {
		config.MackerelService = os.Getenv("MACKEREL_SERVICE")
	}

	if config.MackerelRole == "" {
		config.MackerelRole = os.Getenv("MACKEREL_ROLE")
	}

	if config.MackerelMetric == "" {
		config.MackerelMetric = os.Getenv("MACKEREL_METRIC")
	}

	if config.AWSProfile == "" {
		config.AWSProfile = os.Getenv("AWS_PROFILE")
	}

	if config.AWSRegion == "" {
		config.AWSRegion = os.Getenv("AWS_REGION")
	}

	if config.MackerelAPIKey == "" {
		config.MackerelAPIKey = os.Getenv("MACKEREL_APIKEY")
	}

	// 必須パラメータのバリデーション
	if config.MackerelService == "" {
		fmt.Fprintf(os.Stderr, "Error: Mackerel service is required\n")
		flag.Usage()
		os.Exit(1)
	}

	if config.MackerelRole == "" {
		fmt.Fprintf(os.Stderr, "Error: Mackerel role is required\n")
		flag.Usage()
		os.Exit(1)
	}

	if config.MackerelMetric == "" {
		fmt.Fprintf(os.Stderr, "Error: Mackerel metric name is required\n")
		flag.Usage()
		os.Exit(1)
	}

	if config.MackerelAPIKey == "" {
		fmt.Fprintf(os.Stderr, "Error: Mackerel API key is required\n")
		flag.Usage()
		os.Exit(1)
	}

	return config
}

func maskAPIKey(apiKey string) string {
	if len(apiKey) <= 8 {
		return "***"
	}
	return apiKey[:4] + "****" + apiKey[len(apiKey)-4:]
}
