FROM golang:1.11
WORKDIR $GOPATH/src/github.com/kacejot/rep-controller
Add . .
RUN go get -v -d ./... && go build -v
ENTRYPOINT ["./rep-controller"]
