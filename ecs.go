package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

func StopEcsTask(ecsClient *ecs.Client, task *Task, dryRun bool) error {
	if dryRun {
		log.Printf("Dry run: would stop task with ARN %s", task.TaskArn)
		return nil
	}

	input := &ecs.StopTaskInput{
		Task:    &task.TaskArn,
		Cluster: &task.ClusterArn,
	}

	_, err := ecsClient.StopTask(context.Background(), input)
	if err != nil {
		return fmt.Errorf("failed to stop task %s: %w", task.TaskArn, err)
	}

	log.Printf("Stopped task with ARN %s", task.TaskArn)
	return nil
}
