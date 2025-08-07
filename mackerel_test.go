package main

import (
	"math"
	"testing"
)

func TestFindMostMemoryConsumingTaskFromData(t *testing.T) {
	tests := []struct {
		name     string
		hosts    []HostData
		metrics  []MetricData
		expected *Task
	}{
		{
			name:     "empty hosts and metrics",
			hosts:    []HostData{},
			metrics:  []MetricData{},
			expected: nil,
		},
		{
			name:  "empty hosts",
			hosts: []HostData{},
			metrics: []MetricData{
				{HostID: "host1", Value: 10.0},
			},
			expected: nil,
		},
		{
			name: "empty metrics",
			hosts: []HostData{
				{ID: "host1", TaskArn: "task1", ClusterArn: "cluster1"},
			},
			metrics:  []MetricData{},
			expected: nil,
		},
		{
			name: "single host and metric",
			hosts: []HostData{
				{ID: "host1", TaskArn: "task1", ClusterArn: "cluster1"},
			},
			metrics: []MetricData{
				{HostID: "host1", Value: 10.0},
			},
			expected: &Task{
				TaskArn:    "task1",
				ClusterArn: "cluster1",
			},
		},
		{
			name: "multiple hosts with different memory consumption",
			hosts: []HostData{
				{ID: "host1", TaskArn: "task1", ClusterArn: "cluster1"},
				{ID: "host2", TaskArn: "task2", ClusterArn: "cluster2"},
				{ID: "host3", TaskArn: "task3", ClusterArn: "cluster3"},
			},
			metrics: []MetricData{
				{HostID: "host1", Value: 10.0},
				{HostID: "host2", Value: 25.0}, // highest
				{HostID: "host3", Value: 15.0},
			},
			expected: &Task{
				TaskArn:    "task2",
				ClusterArn: "cluster2",
			},
		},
		{
			name: "host with metric not in hosts list",
			hosts: []HostData{
				{ID: "host1", TaskArn: "task1", ClusterArn: "cluster1"},
			},
			metrics: []MetricData{
				{HostID: "host_unknown", Value: 10.0},
			},
			expected: nil,
		},
		{
			name: "multiple metrics with same host ID (should use highest)",
			hosts: []HostData{
				{ID: "host1", TaskArn: "task1", ClusterArn: "cluster1"},
				{ID: "host2", TaskArn: "task2", ClusterArn: "cluster2"},
			},
			metrics: []MetricData{
				{HostID: "host1", Value: 10.0},
				{HostID: "host1", Value: 5.0},  // lower value for same host
				{HostID: "host2", Value: 15.0}, // highest overall
			},
			expected: &Task{
				TaskArn:    "task2",
				ClusterArn: "cluster2",
			},
		},
		{
			name: "NaN values should be handled correctly",
			hosts: []HostData{
				{ID: "host1", TaskArn: "task1", ClusterArn: "cluster1"},
				{ID: "host2", TaskArn: "task2", ClusterArn: "cluster2"},
			},
			metrics: []MetricData{
				{HostID: "host1", Value: math.NaN()},
				{HostID: "host2", Value: 10.0},
			},
			expected: &Task{
				TaskArn:    "task2",
				ClusterArn: "cluster2",
			},
		},
		{
			name: "all NaN values",
			hosts: []HostData{
				{ID: "host1", TaskArn: "task1", ClusterArn: "cluster1"},
			},
			metrics: []MetricData{
				{HostID: "host1", Value: math.NaN()},
			},
			expected: &Task{
				TaskArn:    "task1",
				ClusterArn: "cluster1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FindMostMemoryConsumingTaskFromData(tt.hosts, tt.metrics)

			if tt.expected == nil {
				if result != nil {
					t.Errorf("expected nil, got %+v", result)
				}
				return
			}

			if result == nil {
				t.Errorf("expected %+v, got nil", tt.expected)
				return
			}

			if result.TaskArn != tt.expected.TaskArn || result.ClusterArn != tt.expected.ClusterArn {
				t.Errorf("expected %+v, got %+v", tt.expected, result)
			}
		})
	}
}

func TestFindMostMemoryConsumingTaskFromData_EdgeCases(t *testing.T) {
	t.Run("identical values should return first found", func(t *testing.T) {
		hosts := []HostData{
			{ID: "host1", TaskArn: "task1", ClusterArn: "cluster1"},
			{ID: "host2", TaskArn: "task2", ClusterArn: "cluster2"},
		}
		metrics := []MetricData{
			{HostID: "host1", Value: 10.0},
			{HostID: "host2", Value: 10.0},
		}

		result := FindMostMemoryConsumingTaskFromData(hosts, metrics)
		if result == nil {
			t.Fatal("expected non-nil result")
		}

		// Should return one of the tasks (order depends on iteration)
		validResults := map[string]bool{
			"task1": true,
			"task2": true,
		}
		if !validResults[result.TaskArn] {
			t.Errorf("unexpected task ARN: %s", result.TaskArn)
		}
	})
}
