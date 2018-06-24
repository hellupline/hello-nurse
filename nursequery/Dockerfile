FROM golang:1.10.2

WORKDIR /go/src/github.com/hellupline/hello-nurse/nursequery/
EXPOSE 8080

ADD ./ ./

RUN make build

CMD ["./nursequery"]
