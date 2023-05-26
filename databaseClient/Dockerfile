FROM golang:1.20-alpine AS build-stage

ARG GOPATH
WORKDIR /app

COPY go.sum go.mod ./
COPY ./pkg/viper@v1.15.0 ./pkg/viper@v1.15.0

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o ./databaseClient

FROM alpine:latest AS build-release-stage

WORKDIR /app

COPY --from=build-stage /app/databaseClient /app/config/locations.yaml .

EXPOSE 8080
CMD ["./databaseClient"]