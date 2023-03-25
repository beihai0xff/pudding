#!/bin/bash

set -e -x

gomodifytags -all -add-tags json,yaml,mapstructure -w --quiet\
   -file configs/kafka_config.go \
   -file configs/log_config.go \
   -file configs/mysql_config.go \
   -file configs/pulsar_config.go \
   -file configs/redis_config.go \
   -file configs/server_config.go


