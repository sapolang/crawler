FROM golang:latest

# WORKDIR $GOPATH/src/github.com/sapolang/crawler

# RUN go get -u -v github.com/henrylee2cn/pholcus
# RUN go build main.go
WORKDIR /root
ADD ./main /root
EXPOSE 9090
ENTRYPOINT ["./main"]
