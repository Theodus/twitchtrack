FROM golang:1.6
RUN apt-get update
ADD . /go
WORKDIR /go
RUN go build main.go
CMD /go/main
EXPOSE 80
