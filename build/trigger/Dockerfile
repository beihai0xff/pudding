FROM golang:latest AS builder
LABEL stage=gobuilder

ENV GOPROXY https://goproxy.io,direct
RUN go install github.com/bufbuild/buf/cmd/buf@v1.9.0 && \
    go install github.com/fatih/gomodifytags@latest && \
    go install github.com/golang/mock/mockgen@v1.6.0
WORKDIR /build/app/trigger
COPY . .
RUN go mod download -x
RUN make build
RUN mkdir /build/bin && cp ./build/bin/scheduler /build/bin/
RUN cp ./cmd/scheduler/config.yaml /build/bin/config.yaml


FROM alpine
RUN apk update --no-cache
RUN apk add --no-cache ca-certificates
WORKDIR /app/scheduler
COPY --from=builder /build/bin ./

EXPOSE 8081 50051
ENTRYPOINT ["./scheduler"]

HEALTHCHECK --interval=10s --timeout=3s \
  CMD curl -fs http://127.0.0.1:8081/pudding/scheduler/v1/ping || exit 1