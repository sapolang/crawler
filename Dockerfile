FROM golang:latest as builder

RUN go get -u -v github.com/henrylee2cn/pholcus
WORKDIR $GOPATH/src/crawler
COPY . .
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o /root/app main.go

FROM scratch
WORKDIR /root
COPY --from=builder /root/app ./

EXPOSE 9090
CMD  [ "./app" ]

