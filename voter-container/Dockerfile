FROM golang:alpine AS build

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /voter-api

FROM alpine:latest AS run

WORKDIR /

COPY --from=build /voter-api /voter-api

EXPOSE 1080

CMD ["/voter-api"]
