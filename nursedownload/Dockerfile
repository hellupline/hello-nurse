FROM golang:1.10.2

WORKDIR /go/src/github.com/hellupline/hello-nurse/nursedownload/

ADD ./ ./

RUN go build -o nursedownload ./

CMD ["./nursedownload"]
