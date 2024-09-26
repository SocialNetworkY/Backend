FROM golang:latest AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /server cmd/user/main.go

FROM alpine:latest
COPY --from=builder server .
COPY --from=builder /app/config/user.yaml /config/config.yaml
CMD ["/server"]