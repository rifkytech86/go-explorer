# why put builder for minify size file
#FROM golang:1.19-alpine as builder


############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder


# install git + ssl ca ceritifates
# Git is for required for fetching the depedencies
# ca-cerificates is reqired to call HTTPS endpoints
#RUN apk update  && apk add --no-cache git ca-certificates tzdata mysql-client && update-ca-certificates
RUN apk update && apk add --no-cache git ca-certificates tzdata mysql-client && update-ca-certificates


# start workingdir
WORKDIR $GOPATH/src/github.com/naonweh-studio/bubbme-backend
COPY . .


RUN echo $PWD && ls -la

# Fetch depecencies
# RUN go get -d -v
RUN go mod download
RUN go mod verify

# build go binary
# GOOS platform= linux, windows, arm
RUN CGO_ENABLE=0 GOOS=linux GOARCH=amd64 go  build -ldflags="-w -s" -a -installsuffix cgo -o bubbme-backend cmd/main.go

RUN echo $PWD && ls -la



############################
# STEP 2 build a small image
############################
FROM golang:1.19-alpine

RUN apk update


# import from builder
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd

WORKDIR /app

## Copy the executeble
COPY --from=builder /go/src/github.com/naonweh-studio/bubbme-backend/bubbme-backend /app/
COPY --from=builder /go/src/github.com/naonweh-studio/bubbme-backend/.env /app/.env
#COPY ./db/migrations /app/db/migrations

COPY custom-entrypoint.sh /app/

COPY  entrypoint.sh /app/
COPY  wait-for-it.sh /app/


RUN echo $PWD && ls -la
RUN chmod +x bubbme-backend

ENTRYPOINT ["/app/bubbme-backend"]

CMD ["-conf", "/app/.env"]


# Build Script Dockerfile  docker build --tag bubbme-backend --progress=plain .
# Create Network -> run : docker network create bubbme-network
# Run Docker Container Golang -> run : docker run -d -p 8070:9090 --name bubbme-backend --network bubbme-network bubbme-backend

#docker run -d --name mysql-bubbme --network bubbme-network -e MYSQL_ROOT_PASSWORD=ZXCasdqwe123! -e MYSQL_DATABASE=bubbme -p 3306:3306 mysql

## MAKE Smapp
##RUN go build -o main cmd/main.go
#
