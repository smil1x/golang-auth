FROM golang:1.23.1

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN go mod download
RUN go build -o golang-auth ./cmd/main.go

CMD ["./golang-auth"]
