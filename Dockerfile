FROM ubuntu:16.04
RUN apt-get update
RUN apt-get install git golang nodejs-legacy npm -y
ADD . /go
WORKDIR /go
RUN npm install --global elm@0.16
RUN elm-make script.elm --output script.js --yes
RUN go build main.go
CMD /go/main
EXPOSE 80
