package main

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	DryRun         bool
	MackerelRole   string
	MackerelMetric string
	AWSProfile     string
	AWSRegion      string
	MackerelAPIKey string
}

func main() {
	config := parseFlags()

	fmt.Printf("Configuration:\n")
	fmt.Printf("  DryRun: %t\n", config.DryRun)
	fmt.Printf("  Mackerel Role: %s\n", config.MackerelRole)
	fmt.Printf("  Mackerel Metric: %s\n", config.MackerelMetric)
	fmt.Printf("  AWS Profile: %s\n", config.AWSProfile)
	fmt.Printf("  AWS Region: %s\n", config.AWSRegion)
	fmt.Printf("  Mackerel API Key: %s\n", maskAPIKey(config.MackerelAPIKey))
}

func parseFlags() *Config {
	config := &Config{}

	// コマンドラインフラグの定義
	flag.BoolVar(&config.DryRun, "dry-run", false, "Dry run mode (no actual task termination)")
	flag.StringVar(&config.MackerelRole, "mackerel-role", "", "Mackerel role name")
	flag.StringVar(&config.MackerelMetric, "mackerel-metric", "", "Mackerel metric name for memory consumption")
	flag.StringVar(&config.AWSProfile, "aws-profile", "", "AWS profile name")
	flag.StringVar(&config.AWSRegion, "aws-region", "", "AWS region")
	flag.StringVar(&config.MackerelAPIKey, "mackerel-api-key", "", "Mackerel API key")

	flag.Parse()

	// 環境変数からの設定（コマンドラインオプションが優先）
	if config.DryRun == false {
		if os.Getenv("DRY_RUN") == "true" {
			config.DryRun = true
		}
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
		config.MackerelAPIKey = os.Getenv("MACKEREL_API_KEY")
	}

	// 必須パラメータのバリデーション
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

	if config.AWSRegion == "" {
		fmt.Fprintf(os.Stderr, "Error: AWS region is required\n")
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
