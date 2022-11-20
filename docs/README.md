#### Docker 编译

##### 编译 Scheduler

```bash
git clone git@github.com:beihai0xff/pudding.git
cd pudding && make docker-build
```





#### Docker 启动

- 运行 Pulsar 作为实时消息队列：

  ```bash
  docker run -d -it \
      -p 6650:6650 \
      -p 8080:8080 \
      -v pulsardata:/pulsar/data \
      -v pulsarconf:/pulsar/conf \
      --name pulsar-standalone \
      apachepulsar/pulsar:2.10.2 \
      bin/pulsar standalone
  
  
  docker run --name mysql -p 3306:3306  -e MYSQL_ROOT_PASSWORD=my-secret-pw -d mysql:8.0.31

- 运行 Redis 作为 broker：

  ```bash
  docker run --name some-redis -d -it -p 6379:6379 redis:latest
  ```

- 运行 Scheduler 模块：

  ```bash
  docker run --name scheduler -d -it -p 8081:8081 --link pulsar-standalone:localhost --link some-redis:localhost scheduler:alpha-1
  ```

  

-   