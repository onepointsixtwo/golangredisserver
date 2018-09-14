FROM golang:1.10.3-alpine
ADD . /go/src/github.com/onepointsixtwo/golangredisserver
WORKDIR /go/src/github.com/onepointsixtwo/golangredisserver
EXPOSE 6379

CMD ["go", "run", "golangredisserver/main.go"]