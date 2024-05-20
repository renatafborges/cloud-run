FROM golang:1.21 as build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -C ./cmd/tempsystem -o cloudrun

FROM scratch
WORKDIR /app
COPY --from=build /app/cmd/tempsystem/.env /app/cmd/tempsystem/cloudrun ./
ENTRYPOINT ["./cloudrun"]
