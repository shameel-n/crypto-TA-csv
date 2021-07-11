FROM golang:1.12.0-alpine3.9
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN apk add git
RUN go get -d ./...
RUN go build -o main .

CMD ["/app/main"]