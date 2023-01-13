#!/bin/bash

# log config
gomodifytags -file configs/log_config.go -struct -all -add-tags json,yaml,mapstructure -w

gomodifytags -file configs/mysql_config.go -struct MySQLConfig -add-tags json,yaml,mapstructure -w

gomodifytags -file configs/pulsar_config.go -all -add-tags json,yaml,mapstructure -w

gomodifytags -file configs/redis_config.go -struct RedisConfig -add-tags json,yaml,mapstructure -w

gomodifytags -file configs/server_config.go -all -add-tags json,yaml,mapstructure -w



