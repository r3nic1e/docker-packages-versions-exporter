FROM golang:1.18 as builder

WORKDIR /usr/src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o docker-package-exporter .

FROM scratch

WORKDIR /usr/app
COPY --from=builder /usr/src/docker-package-exporter .
ENTRYPOINT ["./docker-package-exporter"]