FROM golang:alpine3.16
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.io,direct
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go
CMD ./main