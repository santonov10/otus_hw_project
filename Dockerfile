FROM golang:1.16-buster as otus

WORKDIR /app
RUN GOOS=linux CGO_ENABLED=0

COPY . .
RUN go mod download

RUN go build ./cmd/*

ENTRYPOINT ["/app/imagepreview"]