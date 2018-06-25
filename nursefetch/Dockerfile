FROM golang:1.10.2

WORKDIR /go/src/github.com/hellupline/hello-nurse/nursefetch/

ADD ./ ./

RUN go build -o nursefetch ./

CMD ["./nursefetch"]
