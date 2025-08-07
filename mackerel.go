package main

import (
	"math"

	"github.com/mackerelio/mackerel-client-go"
)

func FindMostMemoryConsumingTaskArn(client *mackerel.Client, service string, role string, metricName string) (string, error) {
	// メモリ使用量が最も多いタスクを検索するロジックを実装
	hosts, err := client.FindHosts(&mackerel.FindHostsParam{
		Service: service,
		Roles:   []string{role},
	})
	if err != nil {
		return "", err
	}
	if len(hosts) == 0 {
		return "", nil // タスクが見つからない場合は空文字を返す
	}

	hostIds := make([]string, len(hosts))
	taskArnByHostId := make(map[string]string, len(hosts))
	for i, host := range hosts {
		hostIds[i] = host.ID
		if meta, ok := host.Meta.Cloud.MetaData.(map[string]interface{}); ok {
			if arn, ok := meta["task_arn"].(string); ok {
				taskArnByHostId[host.ID] = arn
			}
		}
	}

	values, err := client.FetchLatestMetricValues(hostIds, []string{metricName})
	if err != nil {
		return "", err
	}

	largestValueHostId := ""
	largestValue := math.NaN()
	for hostId, metrics := range values {
		for _, metric := range metrics {
			if metric.Name == metricName {
				// metric.Value が interface なので、 float として取り出す
				if value, ok := metric.Value.(float64); ok {
					// メトリックの値をチェックして、最もメモリ使用量が多いタスクを特定
					if math.IsNaN(largestValue) || value > largestValue {
						largestValue = value
						largestValueHostId = hostId
					}
				}
			}
		}
	}

	return taskArnByHostId[largestValueHostId], nil
}
