# ECS Stop Memory Consuming Task Using Mackerel

A command-line tool that identifies and terminates the ECS task with the highest memory consumption based on Mackerel metrics. This tool is designed for ECS environments using `mackerel-container-agent`.

The primary use case is for ECS Services that gradually accumulate memory usage over time (moderate memory leaks). By running this tool periodically (e.g., daily), you can prevent memory exhaustion and maintain service stability.

## Prerequisites

- The target ECS task must have `mackerel-container-agent` attached as a sidecar container.
- The ECS tasks must be monitored by Mackerel.
- AWS credentials must be configured (via AWS CLI, IAM role, or environment variables).
- Go 1.24+ (for building from source)

## Installation

<!--
### Download Binary

Download the latest release from the [releases page](https://github.com/astj/ecs-stop-memory-consuming-task-using-mackerel/releases).
-->

### Build from Source

```bash
git clone https://github.com/astj/ecs-stop-memory-consuming-task-using-mackerel.git
cd ecs-stop-memory-consuming-task-using-mackerel
go build -o ecs-stop-memory-task
```

## Usage

### Basic Usage

```bash
./ecs-stop-memory-task \
  -mackerel-service "your-service" \
  -mackerel-role "your-role" \
  -mackerel-metric "container.memory.app.usage" \
  -mackerel-api-key "your-api-key"
```

### Dry Run Mode

Test the tool without actually stopping tasks:

```bash
./ecs-stop-memory-task \
  -dry-run \
  -mackerel-service "your-service" \
  -mackerel-role "your-role" \
  -mackerel-metric "container.memory.app.usage" \
  -mackerel-api-key "your-api-key"
```

### Using Environment Variables

```bash
export MACKEREL_SERVICE="your-service"
export MACKEREL_ROLE="your-role"
export MACKEREL_METRIC="container.memory.app.usage"
export MACKEREL_APIKEY="your-api-key"
export DRY_RUN="true"  # Optional: enable dry-run mode

./ecs-stop-memory-task
```

### Command-line Options

| Flag | Environment Variable | Required | Description |
|------|---------------------|----------|-------------|
| `-mackerel-service` | `MACKEREL_SERVICE` | Yes | Mackerel service name |
| `-mackerel-role` | `MACKEREL_ROLE` | Yes | Mackerel role name |
| `-mackerel-metric` | `MACKEREL_METRIC` | Yes | Memory metric name (e.g., `container.memory.app.usage`) |
| `-mackerel-api-key` | `MACKEREL_APIKEY` | Yes | Mackerel API key |
| `-dry-run` | `DRY_RUN` | No | Dry run mode (no actual task termination) |
| `-verbose` | - | No | Enable verbose output |


## Example Cron Job

Set up a daily cron job to automatically manage memory-consuming tasks:

```bash
# Run daily at 3 AM
0 3 * * * /path/to/ecs-stop-memory-task -mackerel-service "web-service" -mackerel-role "app" -mackerel-metric "container.memory.app.usage" -mackerel-api-key "your-key" >> /var/log/ecs-memory-cleanup.log 2>&1
```

## Configuration Examples

### For a web application container

```bash
./ecs-stop-memory-task \
  -mackerel-service "web-service" \
  -mackerel-role "app" \
  -mackerel-metric "container.memory.nginx.usage"
```

### For a background worker

```bash
./ecs-stop-memory-task \
  -mackerel-service "worker-service" \
  -mackerel-role "worker" \
  -mackerel-metric "container.memory.worker.usage"
```

## AWS Permissions

The tool requires the following AWS IAM policy:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "ecs:StopTask"
      ],
      "Resource": "*"
    }
  ]
}
```

## License

MIT License - see the [LICENSE](LICENSE) file for details.
