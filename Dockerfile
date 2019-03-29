FROM golang:1.11
WORKDIR $GOPATH/src/github.com/kacejot/rep-controller
COPY . .
RUN go get -d -v ./... && go build -v
