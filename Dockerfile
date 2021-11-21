FROM golang:1.16

WORKDIR /app

COPY ./ /app

RUN go mod download

RUN go get github.com/santonov10/otus_hw_project

ENTRYPOINT CompileDaemon --build="go build -o main" --command=./main