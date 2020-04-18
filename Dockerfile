FROM golang:latest
MAINTAINER devstackq
RUN mkdir /app
ADD . /app
WORKDIR /app 
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/app
ENTRYPOINT /app
CMD ["go run main"]
