FROM golang:1.24-bookworm AS builder

WORKDIR /app

COPY . .

RUN go build -o /app/ecs-stop-memory-consuming-task-using-mackerel .

FROM debian:bookworm-slim
COPY --from=builder /app/ecs-stop-memory-consuming-task-using-mackerel /bin/ecs-stop-memory-consuming-task-using-mackerel

ENTRYPOINT ["/bin/ecs-stop-memory-consuming-task-using-mackerel"]
