FROM golang:1.20-buster AS builder
LABEL stage=gobuilder

ARG APP="broker"
ENV GOPROXY https://goproxy.cn,direct
WORKDIR /builder
COPY . .
RUN make bootstrap
RUN make build/binary APP=${APP}


FROM ubuntu:22.04

ARG APP="broker"
WORKDIR /pudding/app/${APP}
COPY --from=builder /builder/dist/bin/${APP} ./

EXPOSE 8080 50051
ENTRYPOINT ["./server"]
CMD ["--http_port=8080"]

HEALTHCHECK --interval=10s --timeout=3s \
  CMD curl -fs http://127.0.0.1:8080/healthz || exit 1