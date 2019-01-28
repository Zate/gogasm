FROM golang:latest as builder
WORKDIR /go/src/github.com/zate/gogasm/
#RUN go get -u github.com/golang/dep/cmd/dep
COPY *.go /go/src/github.com/zate/gogasm/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gogasm .

FROM scratch
WORKDIR /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/src/github.com/zate/gogasm/gogasm .
#EXPOSE 2086
CMD ["/gogasm", "-live", "napve" ]
