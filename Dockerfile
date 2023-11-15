FROM golang:1.21.1 as builder
WORKDIR /app

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -mod=vendor -o /app-report /app/cmd/*.go

FROM debian:stable-20231030

WORKDIR /app

COPY --from=builder /app/html /app/html
COPY --from=builder /app-report /app/app-report
ENV HTML_DIR=/app/html

RUN chown -R nobody:nogroup /app && chmod -R 0744 /app

EXPOSE 8080
USER nobody
ENTRYPOINT "/app/app-report"
