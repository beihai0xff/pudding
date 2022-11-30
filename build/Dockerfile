FROM golang:latest AS builder
LABEL stage=gobuilder

ARG app="scheduler"
ENV GOPROXY https://goproxy.cn,direct
RUN go install github.com/bufbuild/buf/cmd/buf@v1.9.0 && \
    go install github.com/fatih/gomodifytags@latest && \
    go install github.com/golang/mock/mockgen@v1.6.0
WORKDIR /build/app/${app}
COPY . .
RUN go mod download -x
RUN make build app=${app}
RUN mkdir /build/bin && cp ./build/bin/server /build/bin/
RUN cp ./cmd/${app}/config.yaml /build/bin/config.yaml


FROM alpine

ARG app=scheduler
RUN apk update --no-cache
RUN apk add --no-cache ca-certificates
WORKDIR /app/${app}
COPY --from=builder /build/bin ./

EXPOSE 8081 50051
ENTRYPOINT ["./server"]

HEALTHCHECK --interval=10s --timeout=3s \
  CMD curl -fs http://127.0.0.1:8081/pudding/${app}/v1/ping || exit 1