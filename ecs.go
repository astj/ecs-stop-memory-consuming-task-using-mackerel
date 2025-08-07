package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

func StopEcsTask(ecsClient *ecs.Client, arn string, dryRun bool) error {
	if dryRun {
		fmt.Printf("Dry run: would stop task with ARN %s\n", arn)
		return nil
	}

	input := &ecs.StopTaskInput{
		Task: &arn,
	}

	_, err := ecsClient.StopTask(context.Background(), input)
	if err != nil {
		return fmt.Errorf("failed to stop task %s: %w", arn, err)
	}

	fmt.Printf("Stopped task with ARN %s\n", arn)
	return nil
}
