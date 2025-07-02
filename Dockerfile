 # Golang backend environment build container
FROM golang:1.19 AS build
ENV DEBIAN_FRONTEND noninteractive

WORKDIR /build
ADD . .

ENV GOOS linux 
ENV GOARCH amd64
RUN go build -o bin/main ./cmd/http/main.go
RUN go build -o bin/chef ./cmd/chef/main.go

# Deploy container
FROM ubuntu AS deploy
ENV DEBIAN_FRONTEND noninteractive

# Copying environment grpc server binary to /usr/src/app
WORKDIR /usr/src/app
COPY --from=build /build/bin/main .
COPY --from=build /build/bin/chef .

EXPOSE 5554

# Start The Project
CMD ["/usr/src/app/main"]
