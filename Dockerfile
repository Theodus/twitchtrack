FROM golang:1.5.3
RUN apt-get update
ADD . /go
WORKDIR /go
RUN go build main.go
CMD /go/main
