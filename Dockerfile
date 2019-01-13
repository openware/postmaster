FROM golang:alpine AS build

ENV GO111MODULE=on

WORKDIR /go/src/app

LABEL maintainer="github@shanaakh.pro"

RUN apk add bash ca-certificates git gcc g++ libc-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /go/bin/pigeon ./cmd/pigeon/main.go

FROM alpine

WORKDIR /app

<<<<<<< HEAD
COPY --from=build /go/bin/pigeon ./pigeon
COPY templates/ templates/
=======
COPY --from=build /go/bin/pigeon /app/pigeon
COPY templates .
>>>>>>> eed5717... Add deployment files

ENTRYPOINT ["./pigeon"]