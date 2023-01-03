FROM golang:latest AS builder
LABEL stage=gobuilder

ARG APP="broker"
ENV GOPROXY https://goproxy.cn,direct
RUN make bootstrap
WORKDIR /build/app/${app}
COPY . .
RUN go mod download -x
RUN make build/binary APP=${APP}
RUN mkdir /build/bin && cp ./build/bin/${APP}/* /build/bin/${APP}/


FROM alpine

ARG APP="broker"
RUN apk update --no-cache
RUN apk add --no-cache ca-certificates
WORKDIR /app/${APP}
COPY --from=builder /build/bin/${APP}/ ./

EXPOSE 8081 50051
ENTRYPOINT ["./server"]

HEALTHCHECK --interval=10s --timeout=3s \
  CMD curl -fs http://127.0.0.1:8081/healthz || exit 1