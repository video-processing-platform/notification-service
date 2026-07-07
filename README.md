# Notification Service

A production-ready notification microservice built with Go and gRPC.

This service is responsible for sending email notifications to users after receiving requests from other services such as the Video Processor Service.

---

## Features

* gRPC Server
* Email Notifications
* SMTP Support
* MailHog Integration
* Mailtrap Integration
* Configurable Mail Provider
* Health Check Endpoints
* Prometheus Metrics
* Structured Logging (Zap)
* Clean Architecture
* Docker Support

---

## Architecture

```text
                 Video Processor Service
                           │
                           │ gRPC
                           ▼
                 Notification Service
                           │
          ┌────────────────┴────────────────┐
          │                                 │
       MailHog                           Mailtrap 

```

---

## Project Structure

```text
notification-service/

├── cmd/
│   └── server/
│
├── config/
│
├── internal/
│   ├── application/
│   ├── bootstrap/
│   ├── domain/
│   ├── infrastructure/
│   │   ├── grpc/
│   │   ├── http/
│   │   ├── logger/
│   │   ├── mail/
│   │   └── metrics/
│   └── interfaces/
│       ├── grpc/
│       └── http/
│
├── proto/
│
├── scripts/
│
├── Dockerfile
├── Makefile
└── go.mod
```

---

## Technologies

* Go
* gRPC
* Protocol Buffers
* SMTP
* MailHog
* Mailtrap
* Prometheus
* Zap Logger
* Docker

---

## gRPC Services

### NotificationService

| RPC                     | Description                                    |
|-------------------------|------------------------------------------------|
| SendEmail               | Send a custom email                            |
| SendProcessingCompleted | Notify user when video processing is completed |

## Testing gRPC

### List Available Services

```bash
grpcurl -plaintext localhost:50051 list
```

Expected output:

```text
grpc.reflection.v1.ServerReflection
grpc.reflection.v1alpha.ServerReflection
notification.NotificationService
```

---

### List Available RPC Methods

```bash
grpcurl -plaintext localhost:50051 list notification.NotificationService
```

Expected output:

```text
notification.NotificationService.SendEmail
notification.NotificationService.SendProcessingCompleted
```

---

### Describe SendEmail

```bash
grpcurl -plaintext localhost:50051 describe notification.NotificationService.SendEmail
```

---

### Describe SendProcessingCompleted

```bash
grpcurl -plaintext localhost:50051 describe notification.NotificationService.SendProcessingCompleted
```

---

### Test SendEmail

```bash
grpcurl -plaintext \
-d '{
  "email":"test@example.com",
  "subject":"Hello",
  "body":"This is a test email."
}' \
localhost:50051 \
notification.NotificationService.SendEmail
```

Expected response:

```json
{
  "success": true,
  "message": "Email sent successfully"
}
```

---

### Test SendProcessingCompleted

```bash
grpcurl -plaintext \
-d '{
  "email":"test@example.com",
  "video_id":"video-123"
}' \
localhost:50051 \
notification.NotificationService.SendProcessingCompleted
```

Expected response:

```json
{
  "success": true,
  "message": "Processing completion notification sent"
}
```

---

### Verify MailHog

Open the MailHog web interface:

```text
http://localhost:8025
```

If MailHog is configured as the active provider, the sent emails should appear in the inbox.

---

### Verify Prometheus Metrics

Open:

```text
http://localhost:8080/metrics
```

Example metrics:

```text
notification_emails_sent_total
notification_emails_failed_total
notification_send_duration_seconds
```

---

### Verify Health Endpoints

Health Check

```bash
curl http://localhost:8080/health
```

---

## HTTP Management Endpoints

| Endpoint   | Description        |
|------------|--------------------|
| `/health`  | Health Check       |
| `/metrics` | Prometheus Metrics |

---

## Metrics

The service exposes Prometheus metrics.

Available metrics include:

* `notification_emails_sent_total`
* `notification_emails_failed_total`
* `notification_send_duration_seconds`

---

## Running the Service

### Install dependencies

```bash
go mod download
```

### Generate gRPC files

```bash
protoc \
  --proto_path=proto \
  --go_out=paths=source_relative:. \
  --go-grpc_out=paths=source_relative:. \
  proto/notification.proto
```

### Run

```bash
go run ./cmd/server
```


---

## Code Quality

Run the linter:

```bash
golangci-lint run
```

Format the code:

```bash
gofmt -w .
```

Run all tests:

```bash
go test ./...
```
