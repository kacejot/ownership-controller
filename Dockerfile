FROM golang:1.11
WORKDIR $GOPATH/src/github.com/kacejot/ownership-controller
Add . .
RUN go get -v -d ./... && go build -v
ENTRYPOINT ["./ownership-controller"]
