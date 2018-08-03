FROM golang:1.10.2 as build

WORKDIR /go/src/github.com/hellupline/hello-nurse/
EXPOSE 8080

ENV USER root

ADD ./ ./

RUN make build

FROM scratch

COPY --from=build /go/src/github.com/hellupline/hello-nurse/hello-nurse hello-nurse

CMD ["./hello-nurse"]
