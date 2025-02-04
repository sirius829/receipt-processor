# Build stage

FROM golang:1.22.11 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -o receipt-processor .

# Final Stage
FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/receipt-processor .

EXPOSE 8080

CMD [ "./receipt-processor" ]