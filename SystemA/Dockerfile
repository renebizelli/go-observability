FROM golang:latest AS builder
WORKDIR /app
COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o systema ./cmd/main.go ./cmd/otel.go

FROM scratch
WORKDIR /app
COPY --from=builder /app/systema .
COPY --from=builder /app/.env .

EXPOSE 8080

CMD [ "./systema" ]
