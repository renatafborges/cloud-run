FROM golang:1.21 as build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o systests

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/exec_automated_test.sh /app/systests ./
ENTRYPOINT [ "sh", "./exec_automated_test.sh" ]
