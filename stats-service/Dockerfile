FROM golang:1.23.5

WORKDIR /app

COPY go.mod ./
COPY ./stats-service ./stats-service
COPY ./common ./common

RUN go mod tidy
RUN go mod download

RUN go build -o stats-service ./stats-service

CMD ["./stats-service/stats-service"]
