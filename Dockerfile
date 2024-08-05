FROM golang:latest

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build -o /Server cmd/app/main.go

CMD [ "/Server" ]