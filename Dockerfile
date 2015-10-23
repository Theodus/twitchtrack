FROM golang:1.5
RUN apt-get update
ADD . /go
WORKDIR /go
RUN go build app.go
CMD /go/app
