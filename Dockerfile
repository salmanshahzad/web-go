FROM golang:1.20.7
WORKDIR /usr/src/app
COPY . .
RUN go build -o bin/web-go .
CMD ["./bin/web-go"]
