package main

import (
	"math"

	"github.com/mackerelio/mackerel-client-go"
)

type Task struct {
	TaskArn    string
	ClusterArn string
}

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

	hostIds := make([]string, len(hosts))
	taskArnByHostId := make(map[string]string, len(hosts))
	clusterArnByHostId := make(map[string]string, len(hosts))
	for i, host := range hosts {
		hostIds[i] = host.ID
		// mackerel-container-agent records task ARN and cluster ARN in host metadata
		if meta, ok := host.Meta.Cloud.MetaData.(map[string]interface{}); ok {
			if arn, ok := meta["task_arn"].(string); ok {
				taskArnByHostId[host.ID] = arn
			}
			if arn, ok := meta["cluster"].(string); ok {
				clusterArnByHostId[host.ID] = arn
			}
		}
	}

	values, err := client.FetchLatestMetricValues(hostIds, []string{metricName})
	if err != nil {
		return nil, err
	}

	largestValueHostId := ""
	largestValue := math.NaN()
	for hostId, metrics := range values {
		for _, metric := range metrics {
			if value, ok := metric.Value.(float64); ok {
				if math.IsNaN(largestValue) || value > largestValue {
					largestValue = value
					largestValueHostId = hostId
				}
			}
		}
	}

	return &Task{
		TaskArn:    taskArnByHostId[largestValueHostId],
		ClusterArn: clusterArnByHostId[largestValueHostId],
	}, nil
}
