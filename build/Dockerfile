FROM golang:1.20-buster AS builder
LABEL stage=gobuilder

WORKDIR /builder
ENV GOPROXY https://goproxy.io,direct
# manage dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download -x
COPY . .
RUN make bootstrap
ARG APP
RUN make build/binary APP=${APP}


FROM ubuntu:22.04
RUN apt update -y
RUN apt install -y curl iputils-ping
ARG APP
WORKDIR /pudding/app/${APP}
COPY --from=builder /builder/dist/bin/${APP} ./

EXPOSE 8080 50051
ENTRYPOINT ["./server"]
CMD ["--http-port=8080 --grpc-port=50051"]