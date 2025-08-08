package main

import (
	"math"

	"github.com/mackerelio/mackerel-client-go"
)

type Task struct {
	TaskArn    string
	ClusterArn string
}

// HostData represents the host information extracted from Mackerel
type HostData struct {
	ID         string
	TaskArn    string
	ClusterArn string
}

// MetricData represents metric values for a host
type MetricData struct {
	HostID string
	Value  float64
}

// FindMostMemoryConsumingTaskFromData finds the task with highest memory consumption from given data
// This function is pure and can be easily unit tested
func FindMostMemoryConsumingTaskFromData(hosts []HostData, metrics []MetricData) *Task {
	if len(hosts) == 0 || len(metrics) == 0 {
		return nil
	}

	// Create map for quick lookup of host data by ID
	hostByID := make(map[string]HostData, len(hosts))
	for _, host := range hosts {
		hostByID[host.ID] = host
	}

	// Find the host with the largest metric value
	largestValueHostID := ""
	largestValue := math.NaN()
	for _, metric := range metrics {
		if math.IsNaN(largestValue) || metric.Value > largestValue {
			largestValue = metric.Value
			largestValueHostID = metric.HostID
		}
	}

	// Return the task for the host with largest value
	if host, exists := hostByID[largestValueHostID]; exists {
		return &Task{
			TaskArn:    host.TaskArn,
			ClusterArn: host.ClusterArn,
		}
	}

	return nil
}

// FindMostMemoryConsumingTask finds the ECS task with highest memory consumption using Mackerel API
func FindMostMemoryConsumingTask(client *mackerel.Client, service string, role string, metricName string) (*Task, error) {
	hosts, err := client.FindHosts(&mackerel.FindHostsParam{
		Service: service,
		Roles:   []string{role},
	})
	if err != nil {
		return nil, err
	}
	if len(hosts) == 0 {
		return nil, nil
	}

	// Extract host data
	hostIds := make([]string, len(hosts))
	hostData := make([]HostData, 0, len(hosts))
	for i, host := range hosts {
		hostIds[i] = host.ID

		data := HostData{ID: host.ID}
		// mackerel-container-agent records task ARN and cluster ARN in host metadata
		if meta, ok := host.Meta.Cloud.MetaData.(map[string]interface{}); ok {
			if arn, ok := meta["task_arn"].(string); ok {
				data.TaskArn = arn
			}
			if arn, ok := meta["cluster"].(string); ok {
				data.ClusterArn = arn
			}
		}

		// Only include hosts with task ARN and cluster ARN
		if data.TaskArn != "" && data.ClusterArn != "" {
			hostData = append(hostData, data)
		}
	}

	// Fetch metric values
	values, err := client.FetchLatestMetricValues(hostIds, []string{metricName})
	if err != nil {
		return nil, err
	}

	// Convert metric values to MetricData slice
	var metricData []MetricData
	for hostId, metrics := range values {
		for _, metric := range metrics {
			if value, ok := metric.Value.(float64); ok {
				metricData = append(metricData, MetricData{
					HostID: hostId,
					Value:  value,
				})
			}
		}
	}

	// Use the pure logic function to find the result
	return FindMostMemoryConsumingTaskFromData(hostData, metricData), nil
}
