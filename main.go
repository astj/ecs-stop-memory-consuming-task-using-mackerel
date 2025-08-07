package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/mackerelio/mackerel-client-go"
)

type Config struct {
	DryRun          bool
	Verbose         bool
	MackerelService string
	MackerelRole    string
	MackerelMetric  string
	MackerelAPIKey  string
}

func main() {
	ctx := context.Background()
	c := parseFlags()

	mackerelClient := mackerel.NewClient(c.MackerelAPIKey)
	mackerelClient.Verbose = c.Verbose

	awsConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	ecsClient := ecs.NewFromConfig(awsConfig)

	arn, err := FindMostMemoryConsumingTaskArn(mackerelClient, c.MackerelService, c.MackerelRole, c.MackerelMetric)
	if err != nil {
		log.Fatalf("Error finding most memory consuming task: %v", err)
	}

	if arn == "" {
		log.Println("No memory consuming task found")
		os.Exit(0)
	}

	log.Println("Most memory consuming task ARN:", arn)

	if err := StopEcsTask(ecsClient, arn, c.DryRun); err != nil {
		log.Fatalf("Error stopping task: %v", err)
	}
	log.Println("Task stopped successfully")
}

func parseFlags() *Config {
	config := &Config{}

	flag.BoolVar(&config.DryRun, "dry-run", false, "Dry run mode (no actual task termination)")
	flag.StringVar(&config.MackerelService, "mackerel-service", "", "Mackerel service name")
	flag.StringVar(&config.MackerelRole, "mackerel-role", "", "Mackerel role name")
	flag.StringVar(&config.MackerelMetric, "mackerel-metric", "", "Mackerel metric name for memory consumption like `container.memory.${target container name}.usage`")
	flag.StringVar(&config.MackerelAPIKey, "mackerel-api-key", "", "Mackerel API key")
	flag.BoolVar(&config.Verbose, "verbose", false, "Enable verbose output")

	flag.Parse()

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

	if config.MackerelAPIKey == "" {
		config.MackerelAPIKey = os.Getenv("MACKEREL_APIKEY")
	}

	if config.MackerelService == "" {
		log.Fatalf("Error: Mackerel service is required")
	}

	if config.MackerelRole == "" {
		log.Fatalf("Error: Mackerel role is required")
	}

	if config.MackerelMetric == "" {
		log.Fatalf("Error: Mackerel metric name is required")
	}

	if config.MackerelAPIKey == "" {
		log.Fatalf("Error: Mackerel API key is required")
	}

	return config
}
