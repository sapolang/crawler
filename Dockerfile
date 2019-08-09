FROM golang:latest

# WORKDIR $GOPATH/src/github.com/sapolang/crawler

# RUN go get -u -v github.com/henrylee2cn/pholcus
# RUN go build main.go

EXPOSE 9090
ENTRYPOINT ["./main"]
