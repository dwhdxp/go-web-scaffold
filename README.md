# Web开发通用脚手架
基于gin+gorm+zap+viper+mysql+redis实现的一个web通用脚手架:
1.gorm负责与mysql交互，也可以采用其他框架，例如sqlx；
2.缓存采用go-redis；
3.日志使用高性能zap日志库；
4.通过Viper来读取配置信息；
## Viper配置
Viper：设置默认值、支持从多种格式配置文件（YAML、JSON、TOML、HCL等）中读取配置信息、还可以实时监控和重新读取配置文件。
### 编写配置文件config.yaml
```yaml
# viper对配置的key是大小写不敏感的，另外key:后一定要有一个空格
name: "web_scaffold"
mode: "dev" # dev or release
port: **** # 指定端口号

log:
  level: "debug"
  filename: "./log/web_scaffold.log" 
  max_size: 200
  max_age: 30
  max_backups: 7

mysql:
  host: "127.0.0.1"
  port: 3306
  user: "root"
  password: "******" # 登录密码
  dbname: "DB"
  max_open_conns: 200
  max_idle_conns: 50

redis:
  host: "127.0.0.1"
  port: 6379
  db: 0
  pool_size: 100
```
### 读取配置文件
```go
// 方式1：直接指定配置文件路径
viper.SetConfigFile("./config/config.yaml")
// 方式2：指定配置文件名、类型、搜索路径
viper.SetConfigName("config")
viper.SetConfigType("yaml")
viper.AddConfigPath(".")
viper.AddConfigPath("./config/")

// 读取配置文件
err = viper.ReadInConfig()
if err != nil {
fmt.Printf("Error reading config file, %s", err)
return err
}
```
### 使用结构体变量保存配置信息
通过定义与配置文件对应的结构体，通过反序列化将配置信息保存到结构体变量中
```go
type Config struct {
	Port int `mapstructure:"port"`
	Name string `mapstructure:"name"`
}

var Conf = new(Config)

// 将配置信息保存到全局变量Conf中
if err := viper.Unmarshal(Conf); err != nil {
	panic(fmt.Errorf("unmarshal conf failed, err:%s \n", err))
}
// 监听配置文件变化
viper.WatchConfig()
viper.OnConfigChange(func(e fsnotify.Event) { // 变更时调用回调函数，实时更新到Conf中
	fmt.Printf("Config file changed: %s", e.Name)
	if err := viper.Unmarshal(Conf); err != nil {
		panic(fmt.Errorf("unmarshal conf failed, err:%s \n", err))
	}	
})
```
## Zap日志
在gin框架中默认的中间件Logger()、Recovery()，可以通过Zap日志库设置中间件来接收gin框架默认输出的日志。
```go
r := gin.New()
r.Use(GinLogger(), GinRecovery())
```
初始化日志库后，通过zap.ReplaceGlobals(lg)替换zap包中全局的logger实例，使得可以更简洁的方式调用
```go
zap.L().Info("info msg")
zap.L().Error("error msg", zap.Error(err))
......
```
## Gorm
Gorm官方文档：https://gorm.io/zh_CN/docs/index.htm
## go-redis
go-redis官方文档：https://redis.uptrace.dev/zh/guide/
