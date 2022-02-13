FROM golang:1.17.5-alpine

WORKDIR /app
COPY . ./
RUN go mod download
RUN go build -o /test-minio-server
EXPOSE 1323
CMD ["/test-minio-server"]