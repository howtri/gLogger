FROM golang:1.22 AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /gLogger ./app/main.go

FROM debian:bookworm-slim

COPY --from=build /gLogger /gLogger

EXPOSE 8081 2112

CMD ["/gLogger"]
