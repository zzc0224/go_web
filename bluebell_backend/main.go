package main

import (
	"bluebell_backend/dao/mysql"
	"bluebell_backend/dao/redis"
	"bluebell_backend/logger"
	"bluebell_backend/pkg/snowflake"
	"bluebell_backend/routers"
	"bluebell_backend/settings"
	"fmt"
	"net/http"
)

// @title 论坛系统
// @version 1.0
// @description 简单的论坛系统
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://swagger.io/supprt
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8081
// @BasePath /api/v1
func main() {
	//var confFile string
	//flag.StringVar(&confFile, "conf", "./conf/config.yaml", "配置文件")
	//flag.Parse()
	// 加载配置
	if err := settings.Init(); err != nil {
		fmt.Printf("load config failed, err:%v\n", err)
		return
	}
	if err := logger.Init(settings.Conf.LogConfig, settings.Conf.Mode); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}
	if err := mysql.Init(settings.Conf.MySQLConfig); err != nil {
		fmt.Printf("init mysql failed, err:%v\n", err)
		return
	}
	defer mysql.Close() // 程序退出关闭数据库连接
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		fmt.Printf("init redis failed, err:%v\n", err)
		return
	}
	defer redis.Close()
	if err := snowflake.Init(1); err != nil {
		fmt.Printf("init snowflake failed, err:%v\n", err)
		return
	}
	// 注册路由
	r := routers.SetupRouter()
	r.StaticFS("/image", http.Dir("./image"))
	err := r.Run(fmt.Sprintf(":%d", settings.Conf.Port))
	if err != nil {
		fmt.Printf("run server failed, err:%v\n", err)
		return
	}
}
