FROM golang:latest AS Builder

# enable Go modules support
ENV GO111MODULE=on
WORKDIR /app

# manage dependencies
COPY go.mod .
COPY go.sum .   
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
#RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main 

FROM  scratch
COPY --from=builder /app/ForumX /app/
EXPOSE 6969
ENTRYPOINT ["/app/ForumX"]
# Copy src code from the host and compile it

# FROM golang:latest 
# RUN mkdir /app 
# ADD . /app/ 
# WORKDIR /app 
# RUN go build -o main . 
# CMD ["/app/main"]
    