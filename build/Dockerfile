FROM golang:1.20.2-alpine AS builder
LABEL stage=gobuilder

RUN apk add --no-cache ca-certificates bash git make
ARG APP="broker"
ENV GOPROXY https://goproxy.cn,direct
WORKDIR /builder
COPY . .
RUN make bootstrap
RUN make gen/proto
RUN make build/binary APP=${APP}


FROM golang:1.20.2-alpine

ARG APP="broker"
RUN apk update --no-cache
RUN apk add --no-cache ca-certificates bash
WORKDIR /pudding/app/${APP}
COPY --from=builder /builder/dist/bin/${APP} ./

EXPOSE 8080 50051
ENTRYPOINT ["./server"]
CMD ["--httpPort=8080"]

HEALTHCHECK --interval=10s --timeout=3s \
  CMD curl -fs http://127.0.0.1:8080/healthz || exit 1