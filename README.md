# Security Report Commenting Tool

This tool allows you to read a security report file, process it, and comment on a GitLab merge request using the GitLab API.

## Features

- Reads a JSON security report.
- Comments on a specific merge request in GitLab with findings.
- Uses GitLab API for authentication.

## Requirements

- Golang 1.18+ (if running via binary)
- Docker (if running via Docker)

## Run via Binary

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/your-repo.git
cd your-repo
```

### 2. Build the Project

```bash
go build -o security-report-tool
```

### 3. Run the Tool

```bash
./security-report-tool -f path/to/report.json \
                       -i <gitlab-project-id> \
                       -u <gitlab-project-url> \
                       -n <gitlab-namespace> \
                       -m <merge-request-id> \
                       -t <gitlab-private-token>
```

### 4. Command-Line Flags

- `-f`: Path to the file containing the security report (e.g., `report.json`)
- `-i`: GitLab project ID (e.g., `$CI_PROJECT_ID`)
- `-u`: GitLab project URL (e.g., `$GITLAB_URL`)
- `-n`: Project namespace (e.g., `$CI_PROJECT_PATH`)
- `-m`: Merge request ID (e.g., `$CI_MERGE_REQUEST_IID`)
- `-t`: GitLab private token (your personal or CI/CD token)

## Run via Docker

### 1. Build Docker Image

```bash
docker build -t security-report-tool .
```

### 2. Run the Docker Container

```bash
docker run --rm \
  -v $(pwd)/path/to/report.json:/app/report.json \
  -e GITLAB_PROJECT_ID=<gitlab-project-id> \
  -e GITLAB_URL=<gitlab-project-url> \
  -e GITLAB_NAMESPACE=<gitlab-namespace> \
  -e MERGE_REQUEST_ID=<merge-request-id> \
  -e GITLAB_TOKEN=<gitlab-private-token> \
  security-report-tool
```

### 3. Environment Variables

- `GITLAB_PROJECT_ID`: GitLab project ID (e.g., `$CI_PROJECT_ID`)
- `GITLAB_URL`: GitLab project URL (e.g., `$GITLAB_URL`)
- `GITLAB_NAMESPACE`: GitLab namespace (e.g., `$CI_PROJECT_PATH`)
- `MERGE_REQUEST_ID`: Merge request ID (e.g., `$CI_MERGE_REQUEST_IID`)
- `GITLAB_TOKEN`: GitLab private token (your personal or CI/CD token)

## Dockerfile Example

```dockerfile
FROM golang:1.18-alpine

WORKDIR /app

COPY . .

RUN go build -o security-report-tool

ENTRYPOINT ["/app/security-report-tool"]
```

## License

MIT License
