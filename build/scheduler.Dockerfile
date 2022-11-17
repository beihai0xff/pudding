FROM golang:latest AS builder
LABEL stage=gobuilder
ENV CGO_ENABLED 0
ENV GOOS linux
ENV GOARCH=amd64
ENV GOPROXY https://goproxy.cn,direct
WORKDIR /build/app/scheduler
COPY . .
RUN go mod download
RUN make build
RUN mkdir /build/bin && cp ./build/bin/scheduler /build/bin/
RUN cp ./cmd/scheduler/config.yaml /build/bin/config.yaml


FROM alpine
RUN apk update --no-cache
RUN apk add --no-cache ca-certificates
WORKDIR /app/scheduler
COPY --from=builder /build/bin ./
CMD ["./scheduler"]