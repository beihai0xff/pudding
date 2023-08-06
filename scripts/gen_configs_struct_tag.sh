#!/bin/bash

set -ex

gomodifytags -all -add-tags json,yaml,mapstructure -w --quiet -file configs/kafka_config.go \
	-file configs/log_config.go \
	-file configs/mysql_config.go \
	-file configs/pulsar_config.go \
	-file configs/redis_config.go \
	-file configs/server_base_config.go \
	-file configs/server_broker_config.go \
	-file configs/server_trigger_config.go
