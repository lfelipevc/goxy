FROM golang:alpine3.12

WORKDIR /go/src/server
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["server"]