FROM golang:1.15.0-alpine3.12 as build

WORKDIR /app


COPY go.mod go.sum ./
# COPY go-control-plane/go.mod go-control-plane/go.sum ./go-control-plane/
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/
# RUN chmod +x main

FROM alpine:3.16
COPY --from=build /app/main /

ENTRYPOINT [ "./main" ]
