FROM golang:1.20.5
WORKDIR /usr/src/app
COPY . .
RUN make build
CMD ["make", "run"]
