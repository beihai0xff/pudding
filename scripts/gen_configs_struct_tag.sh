#!/bin/bash

set -ex

cmd="gomodifytags -all -add-tags json,yaml,mapstructure -w"

${cmd} -file configs/kafka_config.go
${cmd} -file configs/log_config.go
${cmd} -file configs/mysql_config.go
${cmd} -file configs/pulsar_config.go
${cmd} -file configs/redis_config.go
${cmd} -file configs/server_base_config.go
${cmd} -file configs/server_broker_config.go
${cmd} -file configs/server_trigger_config.go
