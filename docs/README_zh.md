#### 编译 Pudding

##### Docker 编译 Scheduler

```bash
git clone git@github.com:beihai0xff/pudding.git
cd pudding && make docker-build
```



#### 运行 Pudding

##### Docker 运行 Scheduler 模块

1. 运行 Pulsar 作为实时消息队列：

  ```bash
  docker run --name pulsar-standalone -d -it \
      -p 6650:6650 \
      -p 8080:8080 \
      -v pulsardata:/pulsar/data \
      -v pulsarconf:/pulsar/conf \
      apachepulsar/pulsar:2.10.2 \
      bin/pulsar standalone

2. 运行 Redis 作为 broker：

   ```bash
   docker run --name some-redis -d -it \
       -p 6379:6379 \
       redis:latest
   ```

3. 运行 Scheduler 模块：

   ```bash
   docker run --name pudding.scheduler -d -it \
       -p 8081:8081 \
       --link pulsar-standalone:localhost \
       --link some-redis:localhost \
       scheduler:alpha-1
   ```

   Scheduler 会默认向地址为`redis://default:default@localhost:6379/11`的 Redis 服务和地址为`pulsar://localhost:6650`的 Pulsar 服务创建连接，如果要以已存在的服务作为依赖组件，可通过自定义启动参数来指定连接地址，例如：

   ```bash
   docker run --name pudding.scheduler -d -it \
       -p 8081:8081 \
       -p 50051:50051 \
       scheduler:alpha-1 \
       -redis=redis://default:default@192.168.10.219:6379/11 \
       -pulsar=pulsar://192.168.10.219:6650
   ```

   ##### Docker 运行 Trigger 模块

   1. 运行 MySQL 数据库：

      ```bas
      docker run --name mysql -p 3306:3306  -e MYSQL_ROOT_PASSWORD=my-secret-pw -d mysql:8.0.31
      ```

      

   2. 