FROM golang:1.20-buster AS builder
LABEL stage=gobuilder

ARG APP="broker"
ENV GOPROXY https://goproxy.cn,direct
WORKDIR /builder
COPY . .
RUN make bootstrap
RUN make build/binary APP=${APP}


FROM ubuntu:22.04
RUN apt update
RUN apt install -y curl
ARG APP="broker"
WORKDIR /pudding/app/${APP}
COPY --from=builder /builder/dist/bin/${APP} ./

EXPOSE 8080 50051
ENTRYPOINT ["./server"]
CMD ["--http-port=8080 --grpc-port=50051"]