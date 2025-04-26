# A prometheus exporter to monitor your DockerHub limits

DockerHub is now rate-limiting image pulls. This is a simple prometheus 
exporter that will help you monitor your DockerHub limits.

## Usage

### Create a config file
Create a `config.yaml` file. Feel free to copy from `config.example.yaml` and modify it to your needs.

### Run the exporter

```bash
docker run --rm \
  -v $PWD/config.yaml:/config.yaml \
  -p 9101:9101 \
  -it ghcr.io/jadolg/dockerhub-pull-limit-exporter:v1.0.0
```

## Available metrics
- The rate limit for DockerHub pulls: `dockerhub_pull_limit_total`
- The remaining DockerHub pulls: `dockerhub_pull_remaining_total`
- The time window in seconds to which the limit applies: `dockerhub_pull_limit_window_seconds`
- The time window in seconds to which the remaining pulls apply: `dockerhub_pull_remaining_window_seconds`
- Exporter errors: `dockerhub_pull_errors_total`
