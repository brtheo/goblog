FROM golang:1.22
ADD main.go /cmd/api/main.go
EXPOSE 8080
WORKDIR /cmd/api
ENTRYPOINT [ "go", "run",  "main.go"]