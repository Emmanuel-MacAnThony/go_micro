FROM golang:1.21-alpine AS builder 
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o authApp ./cmd/api 
RUN chmod +x /app/authApp

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/authApp .
CMD [ "/app/authApp" ]