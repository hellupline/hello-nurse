FROM golang:1.10.2

WORKDIR /go/src/github.com/hello-nurse/nursebadge/
EXPOSE 8080

ADD ./ ./

RUN go get

RUN go build -o ./nursebadge ./

CMD ["./nursebadge"]
