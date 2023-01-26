FROM golang:1.17

LABEL maintainer="lapido@gmail.com"

ENV PORT=8080

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
RUN go build ./main.go

EXPOSE 8080

CMD ["./main"]