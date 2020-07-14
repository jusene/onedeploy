package main

import (
	"deploy/step"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"strings"
)

func init() {
	fmt.Println(`
*** 欢迎使用浙江华为自动部署工具 ***
|___  /  |_  | | | | | | | |/ _ \| |  | |  ___|_   _|
   / /     | | | |_| | | | / /_\ \ |  | | |__   | |
  / /      | | |  _  | | | |  _  | |/\| |  __|  | |
./ /___/\__/ / | | | | |_| | | | \  /\  / |___ _| |_
\_____/\____/  \_| |_/\___/\_| |_/\/  \/\____/ \___/

                          `)
	var choice string
	fmt.Print("Ready(Y/N): ")
	fmt.Scanln(&choice)
	if strings.ToLower(choice) != "y" && strings.ToLower(choice) != "yes" {
		log.Fatal("bye")
	}
}

func main() {
	config := ReadConfig()
	switch jump := config.Get("jump").(string); jump {
	case "INIT":
		goto INIT
	case "HARBOR":
		goto HARBOR
	case "NFS":
		goto NFS
	case "REDIS":
		goto REDIS
	case "MYSQL":
		goto MYSQL
	case "RABBITMQ":
		goto RABBITMQ
	case "KUBEAPP":
		goto KUBEAPP
	case "KUBELAB":
		goto KUBELAB
	default:
		fmt.Println("Let's GO")
	}

	log.Println(`
********** 
检查文件
**********`)
	step.FileCheck(config)

	INIT:
		log.Println(`
********** 
初始化服务器
**********`)
	step.SysInit(config)

	HARBOR:
		log.Println(`
********** 
部署Harbor镜像仓库
**********`)
	step.DeployHarbor(config)

	NFS:
		log.Println(`
********** 
部署NFS服务
**********`)
	step.DeployNFS(config)

	REDIS:
		log.Println(`
********** 
部署REDIS服务
**********`)
	step.DeployRedis(config)

	MYSQL:
		log.Println(`
********** 
部署MYSQL服务
**********`)
	step.DeployMysql(config)

	RABBITMQ:
		log.Println(`
********** 
部署RABBITMQ服务
**********`)
	step.DeployRabbitMQ(config)

	KUBEAPP:
		log.Println(`
********** 
部署应用环境K8S服务
**********`)
	step.DeployK8S(config, "app")

	KUBELAB:
		log.Println(`
********** 
部署应用环境K8S服务
**********`)
	step.DeployK8S(config, "lab")
}

/*
* 读取配置文件
 */
func ReadConfig() *viper.Viper {
	viper.SetConfigType("toml")
	viper.SetConfigName("deploy")
	viper.AddConfigPath(".")
	viper.AddConfigPath("config")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}
	}
	return viper.GetViper()
}