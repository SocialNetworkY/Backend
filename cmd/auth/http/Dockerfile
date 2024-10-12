FROM golang:latest AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY .. .

RUN CGO_ENABLED=0 go build -o /server cmd/auth/http/main.go

FROM alpine:latest
COPY --from=builder server .
CMD ["/server"]