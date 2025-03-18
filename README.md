# Web开发通用脚手架
基于gin+gorm+zap+viper实现的一个web通用脚手架。
1.gorm负责与mysql交互，也可以采用其他框架，例如sqlx；
2.缓存采用go-redis；
3.日志使用高性能zap日志库；
4.通过Viper来读取配置信息；
## Viper配置
Viper：设置默认值、支持从多种格式配置文件（YAML、JSON、TOML、HCL等）中读取配置信息、还可以实时监控和重新读取配置文件。
### 编写配置文件config.yaml
```
name: "web_scaffold"
mode: "dev" # dev or release
port: **** # 指定端口号

log:
  level: "debug"
  filename: "web_demo.log"
  max_size: 200
  max_age: 30
  max_backups: 7

mysql:
  host: "127.0.0.1"
  port: 3306
  user: "root"
  password: "abc123"
  dbname: "db1"
  max_open_conns: 200
  max_idle_conns: 50

redis:
  host: "192.168.238.128"
  port: 6379
  db: 0
```
## Zap日志
## Gorm
## SetUpRouter
