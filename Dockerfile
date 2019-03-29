FROM golang:1.11
WORKDIR $GOPATH/bin
COPY rep-controller ./rep-controller
ENTRYPOINT ["./rep-controller"]
