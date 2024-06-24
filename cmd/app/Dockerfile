FROM golang:1.20-alpine AS builder

WORKDIR /helloworld

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/app ./cmd/app
COPY internal ./internal
COPY pkg ./pkg

RUN CGO_ENABLED=0 GOOS=linux go build -o /helloworld-build ./cmd/app

FROM alpine:3.18

WORKDIR /

COPY --from=builder /helloworld-build ./

EXPOSE 8080

CMD ["./helloworld-build"]
