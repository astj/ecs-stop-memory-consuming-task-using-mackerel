FROM golang:1.24-bookworm AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build -o /app/ecs-stop-memory-consuming-task-using-mackerel .

FROM gcr.io/distroless/static-debian12
COPY --from=builder /app/ecs-stop-memory-consuming-task-using-mackerel /bin/ecs-stop-memory-consuming-task-using-mackerel

ENTRYPOINT ["/bin/ecs-stop-memory-consuming-task-using-mackerel"]
