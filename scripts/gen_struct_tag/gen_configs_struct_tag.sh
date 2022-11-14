#!/bin/sh
go install github.com/fatih/gomodifytags@latest

# log config
gomodifytags -file ../../configs/log_config.go -struct LogConfig -add-tags json,yaml,mapstructure -w
gomodifytags -file ../../configs/log_config.go -struct LogFileConfig -add-tags json,yaml,mapstructure -w

gomodifytags -file ../../configs/mysql_config.go -struct MySQLConfig -add-tags json,yaml,mapstructure -w

gomodifytags -file ../../configs/pulsar_config.go -struct PulsarConfig -add-tags json,yaml,mapstructure -w
gomodifytags -file ../../configs/pulsar_config.go -struct PulsarConfig -add-tags json,yaml,mapstructure -w
gomodifytags -file ../../configs/redis_config.go -struct RedisConfig -add-tags json,yaml,mapstructure -w

gomodifytags -file ../../configs/scheduler_config.go -struct SchedulerConfig -add-tags json,yaml,mapstructure -w


