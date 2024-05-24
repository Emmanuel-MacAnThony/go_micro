FROM golang:1.21-alpine AS builder 
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o listenerApp .
RUN chmod +x /app/listenerApp

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/listenerApp .
CMD [ "/app/listenerApp" ]