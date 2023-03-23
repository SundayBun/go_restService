FROM golang:1.20-alpine as builder
COPY --from=builder ./app ./app

ENV config=docker
WORKDIR /app

RUN go get github.com/githubnemo/CompileDaemon
EXPOSE 5000
EXPOSE 7070
ENTRYPOINT CompileDaemon --build="go build main.go" --command=./main
