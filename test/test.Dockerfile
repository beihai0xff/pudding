FROM golang:latest AS tester
LABEL stage=gotest

ARG APP="broker"
ENV GOPROXY https://goproxy.io,direct
WORKDIR /workspace/pudding