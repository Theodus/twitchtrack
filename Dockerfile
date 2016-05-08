FROM ubuntu:16.04
RUN apt-get update
RUN apt-get install golang nodejs-legacy npm -y
RUN npm install -g elm
ADD . /go
WORKDIR /go
RUN elm-make script.elm --output script.js --yes
RUN go build main.go
CMD /go/main
