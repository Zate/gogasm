FROM golang:latest as builder
WORKDIR /go/src/github.com/zate/gogasm/

COPY *.* /go/src/github.com/zate/gogasm/
COPY frontend /go/src/github.com/zate/gogasm/frontend

RUN export GO111MODULE=on && go mod tidy && go mod vendor && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gogasm .

FROM scratch
WORKDIR /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/src/github.com/zate/gogasm/frontend/* /frontend/
COPY --from=builder /go/src/github.com/zate/gogasm/gogasm .

#EXPOSE 2086
CMD ["/gogasm", "-web", "8081" ]
