#### 编译 Pudding

##### Docker 编译 Broker

```bash
git clone git@github.com:beihai0xff/pudding.git
cd pudding && make build/docker app=broker
```

##### Docker 编译 Trigger

```bash
cd pudding && make build/docker app=trigger
```



#### 运行 Pudding

##### 运行 consul 作为注册中心

   ```bash
   docker run --name=consul1 -d -it \
      -p 8500:8500 \
      --restart=always \
      consul agent \
      --server=true \
      --bootstrap-expect=1 \
      --client=0.0.0.0 -ui
   ```
##### Docker 运行 Broker 模块

1. 运行 Pulsar 作为实时消息队列：

   ```bash
   docker run --name pulsar-standalone -d -it \
      -p 6650:6650 \
      -p 8080:8080 \
      -v pulsardata:/pulsar/data \
      -v pulsarconf:/pulsar/conf \
      --restart=always \
      apachepulsar/pulsar:2.10.2 \
      bin/pulsar standalone
   ```

2. 运行 Redis 作为 Storage：

   ```bash
   docker run --name some-redis -d -it \
       -p 6379:6379 \
       --restart=always \
       redis:latest redis-server \
       --appendonly yes \
       --requirepass "default"
   ```

3. 运行 Broker 模块：

   ```bash
   docker run --name pudding.broker -d -it \
       -p 8081:8081 \
       -p 50051:50051 \
       --link pulsar-standalone:localhost \
       --link some-redis:localhost \
       --link consul1:localhost \ 
       broker:alpha-1
   ```

   Broker 会默认向地址为`redis://default:default@localhost:6379/11`的 Redis 服务和地址为`pulsar://localhost:6650`的 Pulsar 服务创建连接，如果要以已存在的服务作为依赖组件，可通过自定义启动参数来指定连接地址，例如：

   ```bash
   docker run --name pudding.broker -d -it \
       -p 8081:8081 \
       -p 50051:50051 \
       broker:alpha-1 \
       -redis=redis://default:default@192.168.10.219:6379/11 \
       -pulsar=pulsar://192.168.10.219:6650
   ```

   ##### Docker 运行 Trigger 模块

   1. 运行 MySQL 数据库：

      ```bash
      docker run --name mysql -d -it \
      	-p 3306:3306  \
      	-e MYSQL_ROOT_PASSWORD=my-secret-pw \
      	mysql:8.0.31
      ```



2. 