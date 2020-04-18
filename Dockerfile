FROM golang:latest
MAINTAINER devstackq
RUN mkdir /app
ADD . /app
WORKDIR /app 
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
ENTRYPOINT /app
RUN go build -o main
CMD ["/app/main"]
