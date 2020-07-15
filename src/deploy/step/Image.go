package step

import (
	"deploy/utils"
	"github.com/spf13/viper"
	"log"
	"strings"
)

func PullImage(config *viper.Viper) {
	utils.DockerLogin(config.Get("registry.cloud").(map[string]interface{})["username"].(string),
		config.Get("registry.cloud").(map[string]interface{})["password"].(string),
		config.Get("registry.cloud").(map[string]interface{})["domain"].(string))
	cloud := strings.Join([]string{config.Get("registry.cloud").(map[string]interface{})["domain"].(string),
		config.Get("registry.cloud").(map[string]interface{})["group"].(string)}, "/")
	local := strings.Join([]string{config.Get("registry.local").(map[string]interface{})["domain"].(string),
		"library"}, "/")

	log.Printf("拉取镜像%s", cloud+"/"+config.Get("application").(map[string]interface{})["eureka"].(string))
	utils.DockerPull(cloud+"/"+config.Get("application").(map[string]interface{})["eureka"].(string),
		config.GetString("version"))

	log.Printf("拉取镜像%s", cloud+"/"+config.Get("application").(map[string]interface{})["config"].(string))
	utils.DockerPull(cloud+"/"+config.Get("application").(map[string]interface{})["config"].(string),
	config.GetString("version"))

	log.Printf("拉取镜像%s", cloud+"/"+config.Get("application").(map[string]interface{})["zuul"].(string))
	utils.DockerPull(cloud+"/"+config.Get("application").(map[string]interface{})["config"].(string),
		config.GetString("version"))

	for _, app := range config.Get("application").(map[string]interface{})["apps"].(map[string]interface{}) {
		log.Printf("拉取镜像%s", cloud+"/"+app.(string))
		utils.DockerPull(cloud+"/"+app.(string), config.GetString("version"))
	}

	for _, addon := range config.Get("package").(map[string]interface{})["addon"].(map[string]interface{}) {
		log.Printf("拉取镜像%s", cloud+"/"+addon.(string))
		name := strings.Split(addon.(string), ":")[0]
		tag := strings.Split(addon.(string), ":")[1]
		utils.DockerPull(cloud+"/"+name, tag)
	}

	log.Printf("推送镜像%s", local+"/"+config.Get("application").(map[string]interface{})["eureka"].(string))
	utils.DockerTagAndPush(cloud+"/"+config.Get("application").(map[string]interface{})["eureka"].(string),
		local+"/"+config.Get("application").(map[string]interface{})["eureka"].(string),
		config.GetString("version"))

	log.Printf("推送镜像%s", local+"/"+config.Get("application").(map[string]interface{})["config"].(string))
	utils.DockerTagAndPush(cloud+"/"+config.Get("application").(map[string]interface{})["config"].(string),
		local+"/"+config.Get("application").(map[string]interface{})["config"].(string),
		config.GetString("version"))

	log.Printf("推送镜像%s", local+"/"+config.Get("application").(map[string]interface{})["zuul"].(string))
	utils.DockerTagAndPush(cloud+"/"+config.Get("application").(map[string]interface{})["zuul"].(string),
		local+"/"+config.Get("application").(map[string]interface{})["zuul"].(string),
		config.GetString("version"))

	for _, app := range config.Get("application").(map[string]interface{})["apps"].(map[string]interface{}) {
		log.Printf("推送镜像%s", local+"/"+app.(string))
		utils.DockerTagAndPush(cloud+"/"+app.(string), local+"/"+app.(string), config.GetString("version"))
	}

	for _, addon := range config.Get("package").(map[string]interface{})["addon"].(map[string]interface{}) {
		log.Printf("推送镜像%s", local+"/"+addon.(string))
		name := strings.Split(addon.(string), ":")[0]
		tag := strings.Split(addon.(string), ":")[1]
		utils.DockerTagAndPush(cloud+"/"+name, local+"/"+name, tag)
	}
}
