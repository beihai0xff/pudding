## 编译 Pudding

#### Docker 编译 Broker

```bash
git clone git@github.com:beihai0xff/pudding.git
cd pudding && make build/docker app=broker
```

#### Docker 编译 Trigger

```bash
cd pudding && make build/docker app=trigger
```



## 运行 Pudding

#### 生成自签名证书
安装 openssl

```bash
cd deployments/docker-compose/configs
./gen_cert.sh
```
