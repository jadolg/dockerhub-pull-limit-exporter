FROM golang:1.24-alpine AS build

RUN adduser --uid 1000 --disabled-password dockerhub-pull-limit-exporter-user

WORKDIR /app

COPY go.mod go.sum /app/
RUN go mod download

COPY . .

ARG VERSION
ARG COMMIT
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags  \
    "-X main.Version=$VERSION -X main.Commit=$COMMIT -X main.Date=$(date +%Y-%m-%dT%H:%M:%SZ)"

FROM scratch
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app/dockerhub-pull-limit-exporter /dockerhub-pull-limit-exporter
USER dockerhub-pull-limit-exporter-user
CMD ["/dockerhub-pull-limit-exporter"]
